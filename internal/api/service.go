package api

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/cainseing/drop-cli/internal/config"
	"github.com/cainseing/drop-cli/internal/crypto"
	"github.com/cainseing/drop-cli/internal/display"
)

func HandleCreateCommand(input []byte, ttl int, reads int) {
	if ttl > config.MaxTTLMinutes {
		display.PrintError("Error", fmt.Sprintf("TTL exceeds maximum allowed limit (%d days)", config.MaxTTLMinutes/1440), nil)
		return
	}

	if ttl <= 0 {
		display.PrintError("Error", "TTL must be at least 1 minute", nil)
		return
	}

	if len(input) > config.MaxBlobSize {
		display.PrintError("Error", fmt.Sprintf("Payload too large (Max: %dKB)", config.MaxBlobSize/1024), nil)
		return
	}

	ciphertext, key, err := crypto.Encrypt(input)

	if err != nil {
		display.PrintError("Error", "Encryption Error", err)
		return
	}

	encodedBlob := base64.StdEncoding.EncodeToString(ciphertext)

	id, err := postBlob(encodedBlob, ttl, reads)

	fmt.Print("\r\033[K")

	if err != nil {
		display.PrintError("Error", "API Error", err)
		return
	}

	rawToken := fmt.Sprintf("%s.%s.%s", config.ProtocolVersion, id, hex.EncodeToString(key))
	token := "drop_" + base64.RawURLEncoding.EncodeToString([]byte(rawToken))

	display.PrintSuccess("Drop token", "")
	fmt.Printf("\n%s\n\n", display.Secret.Render(token))
}

func HandleGetCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		display.PrintError("Error", "Token provided is not valid", nil)
		return
	}
	parts := strings.Split(string(decoded), ".")

	if len(parts) != 3 {
		display.PrintError("Error", "Token provided is not valid", nil)
		return
	}

	usedProtocol, id, keyHex := parts[0], parts[1], parts[2]

	if config.ProtocolVersion != usedProtocol {
		if config.ProtocolVersion > usedProtocol {
			display.PrintError("Error", "This Drop is incompatible because the sender's version is out of date. Please ask them to update their Drop CLI and generate a new Drop.", nil)
			return
		}
		display.PrintError("Error", "To decrypt this Drop, an update is required. Please install the latest version of the Drop CLI.", nil)
		return
	}

	key, err := hex.DecodeString(keyHex)

	if err != nil {
		display.PrintError("Error", "", err)
		return
	}

	response, err := getBlob(id)
	fmt.Fprintf(os.Stderr, "\r\033[K")

	if err != nil {
		display.PrintError("Error", "", err)
		return
	}

	ciphertext, _ := base64.StdEncoding.DecodeString(response.Blob)
	plaintext, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		display.PrintError("Error", "", err)
		return
	}

	stat, _ := os.Stdout.Stat()
	isTerminal := (stat.Mode() & os.ModeCharDevice) != 0

	if !isTerminal {
		fmt.Fprint(os.Stdout, string(plaintext))
		return
	}

	display.PrintSuccess("Drop Decrypted", "")
	fmt.Printf("\n%s\n", display.Secret.Render(string(plaintext)))

	if response.RemainingReads > 0 {
		label := "Reads"
		if response.RemainingReads == 1 {
			label = "Read"
		}

		fmt.Printf("\n\n%s %s\n",
			display.ErrorLabel.Render(fmt.Sprintf("%d", response.RemainingReads)),
			display.ErrorText.Render(fmt.Sprintf("%s Remaining", label)))
	}
}

func HandlePurgeCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)

	if err != nil {
		display.PrintError("Error", "Token provided is not valid", nil)
		return
	}

	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		display.PrintError("Error", "Token provided is not valid", nil)
		return
	}

	id := parts[1]

	result, err := purgeBlob(id)

	if err != nil {
		display.PrintError("Error", "Purge failed", err)
		return
	}

	if !result {
		display.PrintError("Error", "Purge failed", nil)
		return
	}

	display.PrintSuccess("Drop Purged", "")
}
