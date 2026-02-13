package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const MIN_SIZE = 128

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, sealed := ciphertext[:nonceSize], ciphertext[nonceSize:]
	envelope, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, err
	}

	if len(envelope) < 4 {
		return nil, fmt.Errorf("Invalid envolope")
	}

	actualLen := binary.BigEndian.Uint32(envelope[:4])
	if int(actualLen)+4 > len(envelope) {
		return nil, fmt.Errorf("Size mismatch")
	}

	return envelope[4 : 4+actualLen], nil
}

func encrypt(plaintext []byte) ([]byte, []byte, error) {
	// Envelope: [Length][Data][Padding]
	actualLen := uint32(len(plaintext))
	envelope := make([]byte, 4)
	binary.BigEndian.PutUint32(envelope, actualLen)
	envelope = append(envelope, plaintext...)

	if len(envelope) < MIN_SIZE {
		padding := make([]byte, MIN_SIZE-len(envelope))
		rand.Read(padding)
		envelope = append(envelope, padding...)
	}

	// AES-GCM Setup
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, envelope, nil)
	return ciphertext, key, nil
}
