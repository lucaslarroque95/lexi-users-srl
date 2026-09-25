package testutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"lexi-users-srl/internal/lib/utils"
)

// NewKeys builds a utils.Keys backed by a freshly generated RSA key pair,
// so tests can sign/verify real tokens without touching the repo's key files.
func NewKeys(t *testing.T) utils.Keys {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	return NewKeysFromRSA(t, priv)
}

// NewKeysFromRSA builds a utils.Keys around a caller-supplied RSA key pair.
// Useful when a test needs the raw *rsa.PrivateKey (e.g. to hand-craft a
// token with jwt-go directly) while still verifying through utils.Keys.
func NewKeysFromRSA(t *testing.T, priv *rsa.PrivateKey) utils.Keys {
	t.Helper()

	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.pem")
	pubPath := filepath.Join(dir, "public.pem")

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	if err := os.WriteFile(privPath, privPEM, 0o600); err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	if err := os.WriteFile(pubPath, pubPEM, 0o600); err != nil {
		t.Fatalf("failed to write public key: %v", err)
	}

	keys := utils.Keys{}
	if err := keys.LoadPrivateKey(privPath); err != nil {
		t.Fatalf("failed to load generated private key: %v", err)
	}
	if err := keys.LoadPublicKey(pubPath); err != nil {
		t.Fatalf("failed to load generated public key: %v", err)
	}

	return keys
}
