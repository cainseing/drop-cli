package signer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cainseing/drop-cli/internal/identity"
	"golang.org/x/crypto/ssh"
)

type Signer struct{}

func New() *Signer {
	return &Signer{}
}

func (s *Signer) Sign(payload []byte, authorizedKeys []ssh.PublicKey) ([]byte, error) {
	return s.signWithKey(payload, authorizedKeys)
}

func VerifySignature(blobB64, signatureB64, sender, provider string) error {
	authKeys, err := identity.FetchKeys(identity.Provider(provider), sender)
	if err != nil {
		return fmt.Errorf("could not retrieve public keys for %s: %w", sender, err)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return fmt.Errorf("invalid signature encoding")
	}

	actualDataToVerify, err := base64.StdEncoding.DecodeString(blobB64)
	if err != nil {
		return fmt.Errorf("invalid blob encoding")
	}

	for _, pubKey := range authKeys {
		err := pubKey.Verify(actualDataToVerify, &ssh.Signature{
			Format: pubKey.Type(),
			Blob:   sigBytes,
		})

		if err == nil {
			// Match found
			return nil
		}
	}

	return fmt.Errorf("no matching key found for signature")
}

func (s *Signer) signWithKey(payload []byte, authorizedKeys []ssh.PublicKey) ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	paths := []string{
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(home, ".ssh", "id_rsa"),
		filepath.Join(home, ".ssh", "id_ecdsa"),
	}

	for _, path := range paths {
		keyBytes, err := os.ReadFile(path) // #nosec G304 -- path is constructed from hardcoded SSH key names under home dir
		if err != nil {
			continue
		}

		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			continue
		}

		if !s.isKeyAuthorized(signer.PublicKey(), authorizedKeys) {
			continue
		}

		sig, err := signer.Sign(nil, payload)
		if err != nil {
			return nil, err
		}

		return sig.Blob, nil
	}

	return nil, fmt.Errorf("no matching local files")
}

func (s *Signer) isKeyAuthorized(key ssh.PublicKey, authorizedKeys []ssh.PublicKey) bool {
	target := key.Marshal()
	for _, authKey := range authorizedKeys {
		if bytes.Equal(target, authKey.Marshal()) {
			return true
		}
	}
	return false
}
