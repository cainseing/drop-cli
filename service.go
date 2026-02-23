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
		PrintError("Error", fmt.Sprintf("TTL exceeds maximum allowed limit (%d days)", MaxTTLMinutes/1440), nil)
		return
	}

	if ttl <= 0 {
		PrintError("Error", "TTL must be at least 1 minute", nil)
		return
	}

	if len(input) > MaxBlobSize {
		PrintError("Error", fmt.Sprintf("Payload too large (Max: %dKB)", MaxBlobSize/1024), nil)
		return
	}

	ciphertext, key, err := encrypt(input)

	if err != nil {
		PrintError("Error", "Encryption Error", err)
		return
	}

	encodedBlob := base64.StdEncoding.EncodeToString(ciphertext)

	id, err := postBlob(encodedBlob, ttl, reads)

	fmt.Print("\r\033[K")

	if err != nil {
		PrintError("Error", "API Error", err)
		return
	}

	rawToken := fmt.Sprintf("%s.%s.%s", protocolVersion, id, hex.EncodeToString(key))
	token := "drop_" + base64.RawURLEncoding.EncodeToString([]byte(rawToken))

	PrintSuccess("Drop token", "")
	fmt.Printf("\n%s\n\n", secret.Render(token))
}

func handleGetCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		PrintError("Error", "Token provided is not valid", nil)
		return
	}
	parts := strings.Split(string(decoded), ".")

	if len(parts) != 3 {
		PrintError("Error", "Token provided is not valid", nil)
		return
	}

	usedProtocol, id, keyHex := parts[0], parts[1], parts[2]

	if protocolVersion != usedProtocol {
		if protocolVersion > usedProtocol {
			PrintError("Error", "This Drop is incompatible because the sender's version is out of date. Please ask them to update their Drop CLI and generate a new Drop.", nil)
			return
		}
		PrintError("Error", "To decrypt this Drop, an update is required. Please install the latest version of the Drop CLI.", nil)
		return
	}

	key, err := hex.DecodeString(keyHex)

	if err != nil {
		PrintError("Error", "", err)
		return
	}

	response, err := getBlob(id)
	fmt.Fprintf(os.Stderr, "\r\033[K")

	if err != nil {
		PrintError("Error", "", err)
		return
	}

	ciphertext, _ := base64.StdEncoding.DecodeString(response.Blob)
	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		PrintError("Error", "", err)
		return
	}

	stat, _ := os.Stdout.Stat()
	isTerminal := (stat.Mode() & os.ModeCharDevice) != 0

	if !isTerminal {
		fmt.Fprint(os.Stdout, string(plaintext))
		return
	}

	PrintSuccess("Drop Decrypted", "")
	fmt.Printf("\n%s\n", secret.Render(string(plaintext)))

	if response.RemainingReads > 0 {
		label := "Reads"
		if response.RemainingReads == 1 {
			label = "Read"
		}

		fmt.Printf("\n\n%s %s\n",
			errorLabel.Render(fmt.Sprintf("%d", response.RemainingReads)),
			errorText.Render(fmt.Sprintf("%s Remaining", label)))
	}
}

func handlePurgeCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)

	if err != nil {
		PrintError("Error", "Token provided is not valid", nil)
		return
	}

	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		PrintError("Error", "Token provided is not valid", nil)
		return
	}

	id := parts[1]

	result, err := purgeBlob(id)

	if err != nil {
		PrintError("Error", "Purge failed", err)
		return
	}

	if !result {
		PrintError("Error", "Purge failed", nil)
		return
	}

	PrintSuccess("Drop Purged", "")
}
