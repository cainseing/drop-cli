package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	plaintext := []byte("Hello, World!")

	ciphertext, key, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if key == nil || len(key) != 32 {
		t.Fatalf("Expected 32-byte key, got %v bytes", len(key))
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypted plaintext mismatch. Expected %v, got %v", plaintext, decrypted)
	}
}

func TestEncryptSmallPayload(t *testing.T) {
	// Payload smaller than MIN_SIZE should be padded
	plaintext := []byte("tiny")

	ciphertext, key, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Small payload mismatch. Expected %v, got %v", plaintext, decrypted)
	}
}

func TestEncryptLargePayload(t *testing.T) {
	// Payload larger than MIN_SIZE should not be padded
	plaintext := make([]byte, 512)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}

	ciphertext, key, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Large payload mismatch")
	}
}

func TestEncryptEmptyPayload(t *testing.T) {
	plaintext := []byte{}

	ciphertext, key, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Empty payload mismatch. Expected empty, got %v", decrypted)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	plaintext := []byte("Secret data")

	ciphertext, _, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Use a different key
	wrongKey := make([]byte, 32)
	for i := range wrongKey {
		wrongKey[i] = 0xFF
	}

	_, err = Decrypt(ciphertext, wrongKey)
	if err == nil {
		t.Fatalf("Decrypt with wrong key should have failed")
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	invalidCiphertext := []byte("not valid ciphertext")
	key := make([]byte, 32)

	_, err := Decrypt(invalidCiphertext, key)
	if err == nil {
		t.Fatalf("Decrypt with invalid ciphertext should have failed")
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	plaintext := []byte("Authentic message")

	ciphertext, key, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Tamper with the ciphertext
	if len(ciphertext) > 20 {
		ciphertext[20] ^= 0xFF
	}

	_, err = Decrypt(ciphertext, key)
	if err == nil {
		t.Fatalf("Decrypt with tampered ciphertext should have failed")
	}
}

func TestDecryptTooShortCiphertext(t *testing.T) {
	// Ciphertext shorter than nonce size should fail
	tooShort := []byte("x")
	key := make([]byte, 32)

	_, err := Decrypt(tooShort, key)
	if err == nil {
		t.Fatalf("Decrypt with too short ciphertext should have failed")
	}
}

func TestDecryptInvalidEnvelope(t *testing.T) {
	plaintext := []byte("test")

	ciphertext, key, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Truncate to create invalid envelope
	if len(ciphertext) > 20 {
		invalidCiphertext := ciphertext[:20]
		_, err = Decrypt(invalidCiphertext, key)
		if err == nil {
			t.Fatalf("Decrypt with truncated ciphertext should have failed")
		}
	}
}

func TestEncryptProducesUniqueNonces(t *testing.T) {
	plaintext := []byte("Same plaintext")

	ciphertext1, key1, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("First encrypt failed: %v", err)
	}

	ciphertext2, key2, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Second encrypt failed: %v", err)
	}

	// Keys should be different (random generation)
	if bytes.Equal(key1, key2) {
		t.Fatalf("Keys should be random and different")
	}

	// Ciphertexts should be different (different nonces)
	if bytes.Equal(ciphertext1, ciphertext2) {
		t.Fatalf("Ciphertexts should differ due to random nonces")
	}

	// But both should decrypt to same plaintext
	decrypted1, _ := Decrypt(ciphertext1, key1)
	decrypted2, _ := Decrypt(ciphertext2, key2)

	if !bytes.Equal(decrypted1, plaintext) || !bytes.Equal(decrypted2, plaintext) {
		t.Fatalf("Both ciphertexts should decrypt to original plaintext")
	}
}

func TestMinSizePadding(t *testing.T) {
	// Payload of 4 bytes (length header) + some data should trigger padding
	plaintext := []byte("small")
	ciphertext, key, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// The envelope should be padded to at least MIN_SIZE
	// Decrypt and verify the original payload is intact
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Payload mismatch after padding. Expected %v, got %v", plaintext, decrypted)
	}
}

func TestInvalidKeySize(t *testing.T) {
	plaintext := []byte("test")
	ciphertext, _, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Try to decrypt with wrong key size
	shortKey := make([]byte, 16)
	_, err = Decrypt(ciphertext, shortKey)
	if err == nil {
		t.Fatalf("Decrypt with invalid key size should fail")
	}
}
