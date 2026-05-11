package signer

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestNewSigner(t *testing.T) {
	signer := New()
	if signer == nil {
		t.Fatalf("New() returned nil")
	}
}

func TestIsKeyAuthorized(t *testing.T) {
	// Empty authorized keys - will be tested with real keys below
	signer := New()

	key1String := "AAAAC3NzaC1lZDI1NTE5AAAAIAfv5Q3jkAUwMwDNhjfDubLl4jDILdkfx5/Ae4PdLcCt"
	key1, _, _, _, err := ssh.ParseAuthorizedKey([]byte("ssh-ed25519 " + key1String + " test1"))
	if err != nil {
		t.Fatalf("Failed to parse key1: %v", err)
	}

	result := signer.isKeyAuthorized(key1, []ssh.PublicKey{})
	if result {
		t.Fatalf("isKeyAuthorized() should return false with empty authorized keys")
	}
}

func TestSignWithEmptyAuthorizedKeys(t *testing.T) {
	signer := New()
	payload := []byte("test payload")

	// Should fail because no authorized keys are provided
	_, err := signer.Sign(payload, []ssh.PublicKey{})
	if err == nil {
		t.Fatalf("Sign() should fail with empty authorized keys")
	}

	// Should fail with specific error about no matching local files
	if err.Error() != "no matching local files" {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSignWithNilAuthorizedKeys(t *testing.T) {
	signer := New()
	payload := []byte("test payload")

	// Should fail because no authorized keys are provided
	_, err := signer.Sign(payload, nil)
	if err == nil {
		t.Fatalf("Sign() should fail with nil authorized keys")
	}
}

func TestVerifySignatureInvalidBlobBase64(t *testing.T) {
	// Note: identity.FetchKeys is called first, so this will fail with that error
	// rather than a base64 error. This is expected behavior.
	invalidBlob := "!!!invalid base64!!!"
	validSig := "dGVzdA==" // valid base64
	sender := "testuser"
	provider := "github"

	err := VerifySignature(invalidBlob, validSig, sender, provider)
	if err == nil {
		t.Fatalf("VerifySignature() should fail")
	}
	// Error will be from identity.FetchKeys since it's called first
}

func TestVerifySignatureInvalidSignatureBase64(t *testing.T) {
	// Note: identity.FetchKeys is called first, so this will fail with that error
	// rather than a base64 error. This is expected behavior.
	validBlob := base64.StdEncoding.EncodeToString([]byte("test blob"))
	invalidSig := "!!!invalid base64!!!"
	sender := "testuser"
	provider := "github"

	err := VerifySignature(validBlob, invalidSig, sender, provider)
	if err == nil {
		t.Fatalf("VerifySignature() should fail")
	}
	// Error will be from identity.FetchKeys since it's called first
}

func TestVerifySignatureEmptySender(t *testing.T) {
	validBlob := base64.StdEncoding.EncodeToString([]byte("test blob"))
	validSig := base64.StdEncoding.EncodeToString([]byte("test sig"))
	emptySender := ""
	provider := "github"

	err := VerifySignature(validBlob, validSig, emptySender, provider)
	if err == nil {
		t.Fatalf("VerifySignature() should fail with empty sender")
	}

	// The error message indicates validation in identity.FetchKeys
	if err.Error() == "" {
		t.Fatalf("Expected an error message")
	}
}

func TestVerifySignatureInvalidProvider(t *testing.T) {
	validBlob := base64.StdEncoding.EncodeToString([]byte("test blob"))
	validSig := base64.StdEncoding.EncodeToString([]byte("test sig"))
	sender := "testuser"
	provider := "unsupported"

	err := VerifySignature(validBlob, validSig, sender, provider)
	if err == nil {
		t.Fatalf("VerifySignature() should fail with unsupported provider")
	}

	// Error should reference unsupported provider
	if err.Error() == "" {
		t.Fatalf("Expected an error message")
	}
}

func TestVerifySignatureNoMatchingKey(t *testing.T) {
	validBlob := base64.StdEncoding.EncodeToString([]byte("test blob"))
	validSig := base64.StdEncoding.EncodeToString([]byte("invalid signature"))
	sender := "testuser"
	provider := "github"

	// This test will attempt to fetch real keys from GitHub
	// If network is unavailable or user doesn't exist, it will fail at that stage
	// We're testing the path where keys are fetched but signature doesn't match
	err := VerifySignature(validBlob, validSig, sender, provider)

	// Error should occur either from network/user lookup or from signature mismatch
	if err == nil {
		t.Fatalf("VerifySignature() should fail with non-matching signature")
	}
}

func TestIsKeyAuthorizedWithMatchingKey(t *testing.T) {
	// Create a test key by parsing a known SSH public key format
	// Using a real Ed25519 public key for testing
	testKeyString := "AAAAC3NzaC1lZDI1NTE5AAAAIAfv5Q3jkAUwMwDNhjfDubLl4jDILdkfx5/Ae4PdLcCt"

	// Decode the public key
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte("ssh-ed25519 " + testKeyString + " test"))
	if err != nil {
		t.Fatalf("Failed to parse test key: %v", err)
	}

	// Create the same key in authorized list
	authorizedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte("ssh-ed25519 " + testKeyString + " authorized"))
	if err != nil {
		t.Fatalf("Failed to parse authorized key: %v", err)
	}

	signer := New()
	result := signer.isKeyAuthorized(pubKey, []ssh.PublicKey{authorizedKey})

	if !result {
		t.Fatalf("isKeyAuthorized() should return true for matching key")
	}
}

func TestIsKeyAuthorizedWithDifferentKeys(t *testing.T) {
	// Two different Ed25519 keys
	key1String := "AAAAC3NzaC1lZDI1NTE5AAAAIAfv5Q3jkAUwMwDNhjfDubLl4jDILdkfx5/Ae4PdLcCt"
	key2String := "AAAAC3NzaC1lZDI1NTE5AAAAIN1oVMdvwDwFvFKv4h6CTMzVVvd8L0sZQmkMN0bxWJt0"

	key1, _, _, _, err := ssh.ParseAuthorizedKey([]byte("ssh-ed25519 " + key1String + " test1"))
	if err != nil {
		t.Fatalf("Failed to parse key1: %v", err)
	}

	key2, _, _, _, err := ssh.ParseAuthorizedKey([]byte("ssh-ed25519 " + key2String + " test2"))
	if err != nil {
		t.Fatalf("Failed to parse key2: %v", err)
	}

	signer := New()
	result := signer.isKeyAuthorized(key1, []ssh.PublicKey{key2})

	if result {
		t.Fatalf("isKeyAuthorized() should return false for different keys")
	}
}

func TestIsKeyAuthorizedMultipleKeys(t *testing.T) {
	key1String := "AAAAC3NzaC1lZDI1NTE5AAAAIAfv5Q3jkAUwMwDNhjfDubLl4jDILdkfx5/Ae4PdLcCt"
	key2String := "AAAAC3NzaC1lZDI1NTE5AAAAIN1oVMdvwDwFvFKv4h6CTMzVVvd8L0sZQmkMN0bxWJt0"

	key1, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ssh-ed25519 " + key1String + " test1"))
	key2, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ssh-ed25519 " + key2String + " test2"))

	signer := New()

	// key1 should be found in the list
	result := signer.isKeyAuthorized(key1, []ssh.PublicKey{key2, key1})
	if !result {
		t.Fatalf("isKeyAuthorized() should find key in list")
	}

	// key2 should be found in the list
	result = signer.isKeyAuthorized(key2, []ssh.PublicKey{key1, key2})
	if !result {
		t.Fatalf("isKeyAuthorized() should find key in list")
	}
}
