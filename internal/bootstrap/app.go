// Package bootstrap wires the real Postgres-backed services once per Lambda
// cold start. Each function's main() calls New() and hands the *App to its
// Handler; tests build an *App by hand around fakes instead (see
// internal/lib/test/testutil), so New() never runs in a test binary.
package bootstrap

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"lexi-users-srl/internal/lib/db"
	"lexi-users-srl/internal/lib/repositories"
	"lexi-users-srl/internal/lib/service"
	"lexi-users-srl/internal/lib/utils"
)

type App struct {
	UserService *service.UserService
	RoleService *service.RoleService
	Keys        utils.Keys
}

// New connects to Postgres (POSTGRES_* env vars, same as the gin server) and
// loads the RSA signing keys. In deployed Lambdas, JWT_PRIVATE_KEY_SECRET_ARN
// and JWT_PUBLIC_KEY_SECRET_ARN point at the AWS Secrets Manager secrets
// Terraform provisions for this service, and the keys are fetched from there
// at cold start. When those env vars aren't set (local dev without AWS
// credentials), New() falls back to reading PEM files from KEYS_DIR
// (defaults to this repo's own keys/ dir).
func New() *App {
	db.InitDB()

	keys := loadKeys()

	roleRepository := repositories.NewPostgresRoleRepository(db.DB)
	userRepository := repositories.NewPostgresUserRepository(db.DB)

	return &App{
		UserService: service.NewUserService(userRepository, roleRepository, keys),
		RoleService: service.NewRoleService(roleRepository),
		Keys:        keys,
	}
}

// loadKeys loads the RSA signing keys from AWS Secrets Manager when
// JWT_PRIVATE_KEY_SECRET_ARN/JWT_PUBLIC_KEY_SECRET_ARN are set, else falls
// back to local KEYS_DIR-based PEM files for local dev.
func loadKeys() utils.Keys {
	privateKeyArn := os.Getenv("JWT_PRIVATE_KEY_SECRET_ARN")
	publicKeyArn := os.Getenv("JWT_PUBLIC_KEY_SECRET_ARN")

	keys := utils.Keys{}

	if privateKeyArn != "" && publicKeyArn != "" {
		ctx := context.Background()

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			panic("failed to load AWS config: " + err.Error())
		}
		client := secretsmanager.NewFromConfig(cfg)

		privateKeyPEM, err := fetchSecret(ctx, client, privateKeyArn)
		if err != nil {
			panic("failed to fetch private key secret: " + err.Error())
		}
		publicKeyPEM, err := fetchSecret(ctx, client, publicKeyArn)
		if err != nil {
			panic("failed to fetch public key secret: " + err.Error())
		}

		if err := keys.LoadPrivateKeyFromPEM(privateKeyPEM); err != nil {
			panic("failed to load private key: " + err.Error())
		}
		if err := keys.LoadPublicKeyFromPEM(publicKeyPEM); err != nil {
			panic("failed to load public key: " + err.Error())
		}

		return keys
	}

	keysDir := os.Getenv("KEYS_DIR")
	if keysDir == "" {
		keysDir = filepath.Join("..", "keys")
	}

	if err := keys.LoadPrivateKey(filepath.Join(keysDir, "private.pem")); err != nil {
		panic("failed to load private key: " + err.Error())
	}
	if err := keys.LoadPublicKey(filepath.Join(keysDir, "public.pem")); err != nil {
		panic("failed to load public key: " + err.Error())
	}

	return keys
}

// fetchSecret retrieves a secret's string value from AWS Secrets Manager.
func fetchSecret(ctx context.Context, client *secretsmanager.Client, arn string) ([]byte, error) {
	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &arn,
	})
	if err != nil {
		return nil, err
	}
	if out.SecretString == nil {
		return nil, errors.New("secret " + arn + " has no SecretString value")
	}
	return []byte(*out.SecretString), nil
}
