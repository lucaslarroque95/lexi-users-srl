package utils

import (
	"crypto/rsa"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Keys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// LoadPrivateKey reads the PEM-encoded RSA private key from path (local
// dev/testing use) and parses it via LoadPrivateKeyFromPEM.
func (k *Keys) LoadPrivateKey(path string) error {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return k.LoadPrivateKeyFromPEM(keyBytes)
}

// LoadPublicKey reads the PEM-encoded RSA public key from path (local
// dev/testing use) and parses it via LoadPublicKeyFromPEM.
func (k *Keys) LoadPublicKey(path string) error {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return k.LoadPublicKeyFromPEM(keyBytes)
}

// LoadPrivateKeyFromPEM parses raw PEM-encoded RSA private key bytes,
// wherever they came from (a file, an AWS Secrets Manager secret, etc.).
func (k *Keys) LoadPrivateKeyFromPEM(pemBytes []byte) error {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return err
	}

	k.privateKey = key
	return nil
}

// LoadPublicKeyFromPEM parses raw PEM-encoded RSA public key bytes,
// wherever they came from (a file, an AWS Secrets Manager secret, etc.).
func (k *Keys) LoadPublicKeyFromPEM(pemBytes []byte) error {
	key, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return err
	}

	k.publicKey = key
	return nil
}

func (k *Keys) GenerateToken(userId string, email string, roles []string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"role":   strings.Join(roles, ","),
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(k.privateKey)
}

func (k *Keys) VerifyToken(token string) (string, []string, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return k.publicKey, nil
	})

	if err != nil {
		return "", nil, errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return "", nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, errors.New("Invalid token claims.")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", nil, errors.New("invalid userId claim")
	}

	roleClaim, _ := claims["role"].(string)
	var roles []string
	if roleClaim != "" {
		roles = strings.Split(roleClaim, ",")
	}

	return userId, roles, nil
}
