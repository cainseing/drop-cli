package signer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/cainseing/drop-cli/internal/identity"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Signer struct {
	agent agent.ExtendedAgent
}

func New() *Signer {
	socket := os.Getenv("SSH_AUTH_SOCK")
	if socket == "" {
		return &Signer{agent: nil}
	}

	conn, err := net.Dial("unix", socket)
	if err != nil {
		return &Signer{agent: nil}
	}

	return &Signer{
		agent: agent.NewClient(conn),
	}
}

func (s *Signer) Sign(payload []byte, authorizedKeys []ssh.PublicKey) ([]byte, error) {
	sig, err := s.signFromFilesystem(payload, authorizedKeys)
	if err == nil {
		return sig, nil
	}

	if s.agent != nil {
		return s.signFromAgent(payload, authorizedKeys)
	}

	return nil, fmt.Errorf("no matching local keys found and ssh-agent is unavailable")
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

func (s *Signer) signFromFilesystem(payload []byte, authorizedKeys []ssh.PublicKey) ([]byte, error) {
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
		keyBytes, err := os.ReadFile(path)
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

func (s *Signer) signFromAgent(payload []byte, authorizedKeys []ssh.PublicKey) ([]byte, error) {
	agentKeys, err := s.agent.List()
	if err != nil {
		return nil, err
	}

	for _, aKey := range agentKeys {
		parsedAgentKey, err := ssh.ParsePublicKey(aKey.Marshal())
		if err != nil {
			continue
		}

		if !s.isKeyAuthorized(parsedAgentKey, authorizedKeys) {
			continue
		}

		sig, err := s.agent.Sign(parsedAgentKey, payload)
		if err != nil {
			continue
		}

		return sig.Blob, nil
	}

	return nil, fmt.Errorf("no matching keys found in ssh-agent")
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
