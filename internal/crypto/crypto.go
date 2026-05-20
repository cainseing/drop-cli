package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const MIN_SIZE = 128

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
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
		return nil, fmt.Errorf("Size mismatch")
	}

	nonce, sealed := ciphertext[:nonceSize], ciphertext[nonceSize:]
	envelope, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, err
	}

	if len(envelope) < 4 {
		return nil, fmt.Errorf("Invalid envelope")
	}

	actualLen := binary.BigEndian.Uint32(envelope[:4])
	if int(actualLen)+4 > len(envelope) {
		return nil, fmt.Errorf("Size mismatch")
	}

	return envelope[4 : 4+actualLen], nil
}

func Encrypt(plaintext []byte) ([]byte, []byte, error) {
	// Envelope: [Length][Data][Padding]
	if len(plaintext) > math.MaxUint32 {
		return nil, nil, fmt.Errorf("plaintext too large")
	}
	actualLen := uint32(len(plaintext)) // #nosec G115 -- bounds checked above
	envelope := make([]byte, 4)
	binary.BigEndian.PutUint32(envelope, actualLen)
	envelope = append(envelope, plaintext...)

	if len(envelope) < MIN_SIZE {
		padding := make([]byte, MIN_SIZE-len(envelope))
		if _, err := io.ReadFull(rand.Reader, padding); err != nil {
			return nil, nil, err
		}
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

	ciphertext := gcm.Seal(nonce, nonce, envelope, nil)
	return ciphertext, key, nil
}
