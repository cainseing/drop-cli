package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func handleCreateCommand(input []byte, ttl int, reads int) {
	if ttl > MaxTTLMinutes {
		printError(fmt.Sprintf("TTL exceeds maximum allowed limit (%d days)", MaxTTLMinutes/1440), nil)
		return
	}

	if ttl <= 0 {
		printError("TTL must be at least 1 minute", nil)
		return
	}

	if len(input) > MaxBlobSize {
		printError(fmt.Sprintf("Payload too large (Max: %dKB)", MaxBlobSize/1024), nil)
		return
	}

	ciphertext, key, err := encrypt(input)

	if err != nil {
		printError("Encryption Error:", err)
		return
	}

	encodedBlob := base64.StdEncoding.EncodeToString(ciphertext)

	id, err := postBlob(encodedBlob, ttl, reads)

	fmt.Print("\r\033[K")

	if err != nil {
		printError("API Error", err)
		return
	}

	rawToken := fmt.Sprintf("%s.%s.%s", protocolVersion, id, hex.EncodeToString(key))
	token := "drop_" + base64.RawURLEncoding.EncodeToString([]byte(rawToken))

	fmt.Printf("\n%s %s\n", highlight.Render(">"), success.Render("DROP CREATED"))

	// 4. Metadata Section
	// fmt.Printf("\n  %s %s\n", dim.Render("TTL:  "), accent.Render(fmt.Sprintf("%d minutes", ttl)))
	// fmt.Printf("  %s %s\n", dim.Render("READS:"), accent.Render(fmt.Sprintf("%d", reads)))

	// fmt.Printf("\n  %s", dim.Render("TOKEN:"))
	fmt.Printf("\n%s\n\n", secret.Render(token))
}

func handleGetCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		printError("Token provided is not valid", nil)
		return
	}
	parts := strings.Split(string(decoded), ".")

	if len(parts) != 3 {
		printError("Token provided is not valid", nil)
		return
	}

	usedProtocol, id, keyHex := parts[0], parts[1], parts[2]

	if protocolVersion != usedProtocol {
		if protocolVersion > usedProtocol {
			printError("The version of the senders client is out of date, they will need to update Drop CLI and re-create the Drop", nil)
			return
		}
		printError("The version of your client is not able to decrypt this Drop, please update Drop CLI", nil)
		return
	}

	key, err := hex.DecodeString(keyHex)

	if err != nil {
		printError("", err)
		return
	}

	response, err := getBlob(id)
	fmt.Fprintf(os.Stderr, "\r\033[K")

	if err != nil {
		printError("", err)
		return
	}

	ciphertext, _ := base64.StdEncoding.DecodeString(response.Blob)
	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		printError("", err)
		return
	}

	stat, _ := os.Stdout.Stat()
	isTerminal := (stat.Mode() & os.ModeCharDevice) != 0

	if !isTerminal {
		fmt.Fprint(os.Stdout, string(plaintext))
		return
	}

	fmt.Printf("\n%s %s\n", highlight.Render(">"), success.Render("DROP RECEIVED"))
	fmt.Printf("\n%s\n\n", secret.Render(string(plaintext)))

	if response.RemainingReads > 0 {
		fmt.Printf("%s %s %s\n",
			dim.Render("──"),
			errorLabel.Render(fmt.Sprintf("%d", response.RemainingReads)),
			errorText.Render("READS REMAINING"))
	} else {
		fmt.Printf("%s %s\n",
			dim.Render("──"),
			errorText.Render("DROP HAS BEEN PURGED"))
	}
	fmt.Println()
}

func handlePurgeCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)

	if err != nil {
		printError("Token provided is not valid", nil)
		return
	}

	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		printError("Token provided is not valid", nil)
		return
	}

	id := parts[1]

	result, err := purgeBlob(id)

	if err != nil {
		printError("Purge failed", err)
		return
	}

	if !result {
		printError("Purge failed", nil)
		return
	}

	fmt.Printf("\n%s %s\n", highlight.Render(">"), success.Render("DROP PURGED"))
}
