package api

import (
	"encoding/base64"
	"testing"

	"github.com/cainseing/drop-cli/internal/config"
	"github.com/cainseing/drop-cli/internal/identity"
	"golang.org/x/crypto/ssh"
)

func TestHandleCreateCommand_InvalidTTL(t *testing.T) {
	// Mock printError
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	HandleCreateCommand([]byte("test"), config.MaxTTLMinutes+1, 1, false, false)

	if len(errors) != 1 || errors[0] != "TTL exceeds maximum allowed limit (7 days)" {
		t.Errorf("Expected TTL error, got %v", errors)
	}
}

func TestHandleCreateCommand_ZeroTTL(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	HandleCreateCommand([]byte("test"), 0, 1, false, false)

	if len(errors) != 1 || errors[0] != "TTL must be at least 1 minute" {
		t.Errorf("Expected TTL zero error, got %v", errors)
	}
}

func TestHandleCreateCommand_PayloadTooLarge(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	largeInput := make([]byte, config.MaxBlobSize+1)
	HandleCreateCommand(largeInput, 5, 1, false, false)

	if len(errors) != 1 || errors[0] != "Payload too large (Max: 1024KB)" {
		t.Errorf("Expected payload size error, got %v", errors)
	}
}

func TestHandleCreateCommand_Success(t *testing.T) {
	// Mock all dependencies
	originalLoadUserConfig := loadUserConfig
	originalEncrypt := encrypt
	originalPostBlobFunc := postBlobFunc
	originalWriteClipboard := writeClipboard
	originalPrintProperty := printProperty

	loadUserConfig = func() *config.UserConfig {
		return &config.UserConfig{Username: "testuser", Provider: "github"}
	}
	encrypt = func(input []byte) ([]byte, []byte, error) {
		return []byte("ciphertext"), []byte("key"), nil
	}
	postBlobFunc = func(blob string, ttl int, reads int, sig string, sender string, provider string) (string, error) {
		return "test-id", nil
	}
	writeClipboard = func(text string) error {
		return nil
	}
	printProperty = func(key, value string) {}

	defer func() {
		loadUserConfig = originalLoadUserConfig
		encrypt = originalEncrypt
		postBlobFunc = originalPostBlobFunc
		writeClipboard = originalWriteClipboard
		printProperty = originalPrintProperty
	}()

	// Capture stdout to check token output
	// For simplicity, just ensure no panic
	HandleCreateCommand([]byte("test"), 5, 1, false, false)
}

func TestHandleGetCommand_InvalidToken(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	HandleGetCommand("invalid")

	if len(errors) != 1 || errors[0] != "Token provided is not valid" {
		t.Errorf("Expected invalid token error, got %v", errors)
	}
}

func TestHandleGetCommand_InvalidParts(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	// Token with wrong number of parts
	HandleGetCommand("drop_" + "invalid")

	if len(errors) != 1 || errors[0] != "Token provided is not valid" {
		t.Errorf("Expected invalid parts error, got %v", errors)
	}
}

func TestHandleGetCommand_SenderProtocolOutOfDate(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	// Token from protocol "1", current config.ProtocolVersion is "2".
	HandleGetCommand("drop_" + "MS5pZC5rZXk6MA") // 1.id.key:0

	want := "This Drop is incompatible because the sender's version is out of date. Please ask them to update their Drop CLI and generate a new Drop."
	if len(errors) != 1 || errors[0] != want {
		t.Errorf("Expected sender-out-of-date error, got %v", errors)
	}
}

func TestHandleGetCommand_RecipientUpdateRequired(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	// Token from protocol "3", current config.ProtocolVersion is "2".
	HandleGetCommand("drop_" + "My5pZC5rZXk6MA") // 3.id.key:0

	want := "To decrypt this Drop, an update is required. Please install the latest version of the Drop CLI."
	if len(errors) != 1 || errors[0] != want {
		t.Errorf("Expected recipient-update-required error, got %v", errors)
	}
}

func TestHandleGetCommand_InvalidKey(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	// Token with invalid hex key
	HandleGetCommand("drop_" + "Mi5pZC5pbnZhbGlkOjA") // 2.id.invalid:0

	if len(errors) != 1 || errors[0] != "Invalid encryption key in token" {
		t.Errorf("Expected invalid key error, got %v", errors)
	}
}

func TestHandleGetCommand_SignedDropMissingSignature(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	originalGetBlobFunc := getBlobFunc
	getBlobFunc = func(id string) (*GetDropResponse, error) {
		return &GetDropResponse{
			Blob: base64.StdEncoding.EncodeToString([]byte("ciphertext")),
			// Signature, Sender, Provider intentionally omitted, simulating
			// a server that stripped them from a drop signed at creation.
		}, nil
	}
	defer func() { getBlobFunc = originalGetBlobFunc }()

	originalDecrypt := decrypt
	decrypt = func(ciphertext []byte, key []byte) ([]byte, error) {
		return []byte("secret"), nil
	}
	defer func() { decrypt = originalDecrypt }()

	// Token with protocol 2, signed_flag = 1, but the server response above
	// has no signature.
	HandleGetCommand("drop_" + "Mi50ZXN0LWlkLjAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA6MQ")

	if len(errors) != 1 || errors[0] != "This drop was signed by the sender but the signature is missing. The content may have been tampered with." {
		t.Errorf("Expected signature-stripped error, got %v", errors)
	}
}

func TestHandlePurgeCommand_InvalidToken(t *testing.T) {
	var errors []string
	originalPrintError := printError
	printError = func(msg string, err error) {
		errors = append(errors, msg)
	}
	defer func() { printError = originalPrintError }()

	HandlePurgeCommand("invalid")

	if len(errors) != 1 || errors[0] != "Token provided is not valid" {
		t.Errorf("Expected invalid token error, got %v", errors)
	}
}

func TestHandlePurgeCommand_Success(t *testing.T) {
	var successes []string
	originalPrintSuccess := printSuccess
	printSuccess = func(msg, detail string) {
		successes = append(successes, msg)
	}
	defer func() { printSuccess = originalPrintSuccess }()

	originalPurgeBlobFunc := purgeBlobFunc
	purgeBlobFunc = func(id string) (bool, error) {
		return true, nil
	}
	defer func() { purgeBlobFunc = originalPurgeBlobFunc }()

	// Valid token: 2.id.key:0
	HandlePurgeCommand("drop_" + "Mi5pZC5rZXk6MA") // 2.id.key:0

	if len(successes) != 1 || successes[0] != "Purged" {
		t.Errorf("Expected success, got %v", successes)
	}
}

// Mock types for signer
type mockSigner struct{}

func (m *mockSigner) Sign(data []byte, keys []ssh.PublicKey) ([]byte, error) {
	return []byte("signature"), nil
}

func TestHandleCreateCommand_Signed(t *testing.T) {
	originalLoadUserConfig := loadUserConfig
	originalEncrypt := encrypt
	originalNewSigner := newSigner
	originalFetchKeys := fetchKeys
	originalPostBlobFunc := postBlobFunc
	originalPrintProperty := printProperty

	loadUserConfig = func() *config.UserConfig {
		return &config.UserConfig{Username: "testuser", Provider: "github"}
	}
	encrypt = func(input []byte) ([]byte, []byte, error) {
		return []byte("ciphertext"), []byte("key"), nil
	}
	newSigner = func() signerInterface {
		return &mockSigner{}
	}
	fetchKeys = func(provider identity.Provider, username string) ([]ssh.PublicKey, error) {
		return []ssh.PublicKey{}, nil
	}
	postBlobFunc = func(blob string, ttl int, reads int, sig string, sender string, provider string) (string, error) {
		return "test-id", nil
	}
	printProperty = func(key, value string) {}

	defer func() {
		loadUserConfig = originalLoadUserConfig
		encrypt = originalEncrypt
		newSigner = originalNewSigner
		fetchKeys = originalFetchKeys
		postBlobFunc = originalPostBlobFunc
		printProperty = originalPrintProperty
	}()

	HandleCreateCommand([]byte("test"), 5, 1, true, false)
}
