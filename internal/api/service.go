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
	"github.com/cainseing/drop-cli/internal/identity"
	"github.com/cainseing/drop-cli/internal/output"
	"github.com/cainseing/drop-cli/internal/signer"
	"golang.org/x/crypto/ssh"
)

var printError = output.PrintError
var loadUserConfig = config.LoadUserConfig
var encrypt = crypto.Encrypt

type signerInterface interface {
	Sign(payload []byte, authorizedKeys []ssh.PublicKey) ([]byte, error)
}

var newSigner = func() signerInterface { return signer.New() }
var fetchKeys = identity.FetchKeys
var postBlobFunc = postBlob
var writeClipboard = clipboard.WriteAll
var stdoutStat = func() (os.FileInfo, error) { return os.Stdout.Stat() }
var printProperty = output.PrintProperty
var printPropertyToStderr = output.PrintPropertyToStderr
var printSuccess = output.PrintSuccess
var verifySignature = signer.VerifySignature
var decrypt = crypto.Decrypt
var getBlobFunc = getBlob
var purgeBlobFunc = purgeBlob

func HandleCreateCommand(input []byte, ttl int, reads int, signed bool, shouldCopy bool) {
	if ttl > config.MaxTTLMinutes {
		printError(fmt.Sprintf("TTL exceeds maximum allowed limit (%d days)", config.MaxTTLMinutes/1440), nil)
		return
	}

	if ttl <= 0 {
		printError("TTL must be at least 1 minute", nil)
		return
	}

	if reads < 1 {
		printError("Reads must be at least 1", nil)
		return
	}

	if len(input) > config.MaxBlobSize {
		printError(fmt.Sprintf("Payload too large (Max: %dKB)", config.MaxBlobSize/1024), nil)
		return
	}

	cfg := loadUserConfig()
	ciphertext, key, err := encrypt(input)
	if err != nil {
		printError("Encryption Error", err)
		return
	}

	var signature []byte
	if signed {
		if cfg.Username == "" || cfg.Provider == "" {
			printError("You need to configure your identity before signing a drop", nil)
			return
		}

		s := newSigner()

		authKeys, err := fetchKeys(identity.Provider(cfg.Provider), cfg.Username)
		if err != nil {
			printError("Could not fetch signing keys", err)
			return
		}

		signature, err = s.Sign(ciphertext, authKeys)
		if err != nil {
			printError("Signing Error", err)
			return
		}
	}

	encodedBlob := base64.StdEncoding.EncodeToString(ciphertext)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	id, err := postBlobFunc(encodedBlob, ttl, reads, encodedSignature, cfg.Username, cfg.Provider)

	if err != nil {
		printError("API Error", err)
		return
	}

	signedFlag := "0"
	if signed {
		signedFlag = "1"
	}

	// The signed flag is packed into the third segment as "hexKey:signedFlag"
	// rather than as a fourth dot-separated segment. This keeps the token at
	// 3 dot-separated parts, so older clients (which check for exactly 3
	// parts before comparing protocol versions) still reach the protocol
	// version check and show a "please update" message instead of "Token
	// provided is not valid".
	rawToken := fmt.Sprintf("%s.%s.%s:%s", config.ProtocolVersion, id, hex.EncodeToString(key), signedFlag)
	token := "drop_" + base64.RawURLEncoding.EncodeToString([]byte(rawToken))

	fmt.Printf("%s\n\n", output.RenderToken(token))

	if signed {
		fmt.Fprintln(os.Stderr, output.DimText.Render(fmt.Sprintf("- Signed via %s by %s", cfg.Provider, cfg.Username)))
	}

	expiryUnit := "minutes"
	if ttl == 1 {
		expiryUnit = "minute"
	}
	readUnit := "Secret will be destroyed on read"
	if reads > 1 {
		readUnit = fmt.Sprintf("%d reads remaining", reads)
	}
	fmt.Fprintln(os.Stderr, output.DimText.Render(fmt.Sprintf("- Expires in %d %s", ttl, expiryUnit)))
	fmt.Fprintln(os.Stderr, output.DimText.Render(fmt.Sprintf("- %s", readUnit)))

	if shouldCopy {
		err := writeClipboard(token)
		if err != nil {
			printError("Could not copy token automatically", err)
		} else {
			fmt.Fprintln(os.Stderr, output.DimText.Render("- Copied to clipboard"))
		}
	}
}

func HandleGetCommand(token string) {
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

	usedProtocol, id, keyAndFlag := parts[0], parts[1], parts[2]

	if config.ProtocolVersion != usedProtocol {
		if config.ProtocolVersion > usedProtocol {
			printError("This Drop is incompatible because the sender's version is out of date. Please ask them to update their Drop CLI and generate a new Drop.", nil)
			return
		}
		printError("To decrypt this Drop, an update is required. Please install the latest version of the Drop CLI.", nil)
		return
	}

	// The third segment is "hexKey:signedFlag" (see HandleCreateCommand).
	keyParts := strings.SplitN(keyAndFlag, ":", 2)
	if len(keyParts) != 2 {
		printError("Token provided is not valid", nil)
		return
	}
	keyHex, signedFlag := keyParts[0], keyParts[1]

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		printError("Invalid encryption key in token", err)
		return
	}

	response, err := getBlobFunc(id)

	if err != nil {
		printError("", err)
		return
	}

	ciphertext, err := base64.StdEncoding.DecodeString(response.Blob)
	if err != nil {
		printError("Failed to decode encrypted payload", err)
		return
	}

	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		printError("", err)
		return
	}

	if signedFlag == "1" {
		fmt.Println()
		if response.Signature == "" || response.Sender == "" || response.Provider == "" {
			printError("This drop was signed by the sender but the signature is missing. The content may have been tampered with.", nil)
			return
		}

		err := verifySignature(response.Blob, response.Signature, response.Sender, response.Provider)
		if err != nil {
			printError("Signature verification failed. The content may have been tampered with.", err)
			return
		}

		fmt.Fprintln(os.Stderr, output.DimText.Render(fmt.Sprintf("- Verified via %s by %s", response.Provider, response.Sender)))
	}

	if response.RemainingReads > 0 {
		fmt.Println()
		label := "reads"
		if response.RemainingReads == 1 {
			label = "read"
		}

		fmt.Fprintln(os.Stderr, output.DimText.Render(fmt.Sprintf("- %d %s remaining", response.RemainingReads, label)))
	}

	fi, err := stdoutStat()
	if err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
		fmt.Fprint(os.Stdout, string(plaintext))
		return
	}

	fmt.Println()
	fmt.Println(output.Secret.Render(string(plaintext)))
}

func HandlePurgeCommand(token string) {
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

	result, err := purgeBlobFunc(id)

	if err != nil {
		printError("", err)
		return
	}

	if !result {
		printError("Unknown error occurred while trying to purge", nil)
		return
	}

	printSuccess("Purged", "")
}
