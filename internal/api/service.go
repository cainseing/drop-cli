package api

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/cainseing/drop-cli/internal/config"
	"github.com/cainseing/drop-cli/internal/crypto"
	"github.com/cainseing/drop-cli/internal/display"
	"github.com/cainseing/drop-cli/internal/identity"
	"github.com/cainseing/drop-cli/internal/signer"
)

func HandleCreateCommand(input []byte, ttl int, reads int, signed bool, shouldCopy bool) {
	if ttl > config.MaxTTLMinutes {
		display.PrintError(fmt.Sprintf("TTL exceeds maximum allowed limit (%d days)", config.MaxTTLMinutes/1440), nil)
		return
	}

	if ttl <= 0 {
		display.PrintError("TTL must be at least 1 minute", nil)
		return
	}

	if len(input) > config.MaxBlobSize {
		display.PrintError(fmt.Sprintf("Payload too large (Max: %dKB)", config.MaxBlobSize/1024), nil)
		return
	}

	cfg := config.LoadUserConfig()
	ciphertext, key, err := crypto.Encrypt(input)
	if err != nil {
		display.PrintError("Encryption Error", err)
		return
	}

	var signature []byte
	if signed {
		if cfg.Username == "" || cfg.Provider == "" {
			display.PrintError("You need to configure your identity before signing a drop", nil)
			return
		}

		s := signer.New()

		authKeys, err := identity.FetchKeys(identity.Provider(cfg.Provider), cfg.Username)
		if err != nil {
			return
		}

		signature, err = s.Sign(ciphertext, authKeys)
		if err != nil {
			display.PrintError("Signing Error", err)
			return
		}
	}

	encodedBlob := base64.StdEncoding.EncodeToString(ciphertext)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	id, err := postBlob(encodedBlob, ttl, reads, encodedSignature, cfg.Username, cfg.Provider)

	fmt.Print("\r\033[K")

	if err != nil {
		display.PrintError("API Error", err)
		return
	}

	rawToken := fmt.Sprintf("%s.%s.%s", config.ProtocolVersion, id, hex.EncodeToString(key))
	token := "drop_" + base64.RawURLEncoding.EncodeToString([]byte(rawToken))

	if signed {
		fmt.Println()
		display.PrintProperty("STATUS", display.StatusVerified.Render(fmt.Sprintf("Signed via %s", cfg.Provider)))
		display.PrintProperty("SENDER", fmt.Sprintf("%s", cfg.Username))
	}

	if shouldCopy {
		err := clipboard.WriteAll(token)
		if err != nil {
			display.PrintError("Could not copy token automatically", err)
		}
	}

	fmt.Printf("%s\n\n", display.Token.Render(token))
}

func HandleGetCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		display.PrintError("Token provided is not valid", nil)
		return
	}
	parts := strings.Split(string(decoded), ".")

	if len(parts) != 3 {
		display.PrintError("Token provided is not valid", nil)
		return
	}

	usedProtocol, id, keyHex := parts[0], parts[1], parts[2]

	if config.ProtocolVersion != usedProtocol {
		if config.ProtocolVersion > usedProtocol {
			display.PrintError("This Drop is incompatible because the sender's version is out of date. Please ask them to update their Drop CLI and generate a new Drop.", nil)
			return
		}
		display.PrintError("To decrypt this Drop, an update is required. Please install the latest version of the Drop CLI.", nil)
		return
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		display.PrintError("Invalid encryption key in token", err)
		return
	}

	response, err := getBlob(id)

	if err != nil {
		display.PrintError("", err)
		return
	}

	ciphertext, err := base64.StdEncoding.DecodeString(response.Blob)
	if err != nil {
		display.PrintError("Failed to decode encrypted payload", err)
		return
	}

	plaintext, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		display.PrintError("", err)
		return
	}

	fi, _ := os.Stdout.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		fmt.Fprint(os.Stdout, string(plaintext))
		return
	}

	fmt.Print("\r\033[K")

	if response.Signature != "" {
		err := signer.VerifySignature(response.Blob, response.Signature, response.Sender, response.Provider)
		if err != nil {
			display.PrintError("Signature verification failed. The content may have been tampered with.", err)
			return
		}

		fmt.Fprintln(os.Stderr)
		display.PrintPropertyToStderr("STATUS", display.StatusVerified.Render(fmt.Sprintf("Verified via %s", response.Provider)))
		display.PrintPropertyToStderr("SENDER", fmt.Sprintf("%s", response.Sender))
	}

	if response.RemainingReads > 0 {
		label := "reads"
		if response.RemainingReads == 1 {
			label = "read"
		}

		fmt.Fprintln(os.Stderr)
		statusText := fmt.Sprintf("%d %s remaining", response.RemainingReads, label)
		display.PrintPropertyToStderr("READS", statusText)
	}

	fmt.Println(display.Secret.Render(string(plaintext)))
}

func HandlePurgeCommand(token string) {
	token = strings.TrimPrefix(token, "drop_")
	decoded, err := base64.RawURLEncoding.DecodeString(token)

	if err != nil {
		display.PrintError("Token provided is not valid", nil)
		return
	}

	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		display.PrintError("Token provided is not valid", nil)
		return
	}

	id := parts[1]

	result, err := purgeBlob(id)

	if err != nil {
		display.PrintError("", err)
		return
	}

	if !result {
		display.PrintError("Unknown error occurred while trying to purge", nil)
		return
	}

	display.PrintSuccess("Purged", "")
}
