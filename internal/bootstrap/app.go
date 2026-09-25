// Package bootstrap wires the real Postgres-backed services once per Lambda
// cold start. Each function's main() calls New() and hands the *App to its
// Handler; tests build an *App by hand around fakes instead (see
// lexi-users-ms/src/test/testutil), so New() never runs in a test binary.
package bootstrap

import (
	"os"
	"path/filepath"

	"lexi/users/db"
	"lexi/users/repositories"
	"lexi/users/service"
	"lexi/users/utils"
)

type App struct {
	UserService *service.UserService
	RoleService *service.RoleService
	Keys        utils.Keys
}

// New connects to Postgres (POSTGRES_* env vars, same as the gin server) and
// loads the RSA keys from KEYS_DIR (defaults to this pilot's own copy of
// lexi-users-ms/src/keys, so it works out of the box before any deploy config
// exists).
func New() *App {
	db.InitDB()

	keysDir := os.Getenv("KEYS_DIR")
	if keysDir == "" {
		keysDir = filepath.Join("..", "..", "microservicios", "lexi-users-ms", "src", "keys")
	}

	keys := utils.Keys{}
	if err := keys.LoadPrivateKey(filepath.Join(keysDir, "private.pem")); err != nil {
		panic("failed to load private key: " + err.Error())
	}
	if err := keys.LoadPublicKey(filepath.Join(keysDir, "public.pem")); err != nil {
		panic("failed to load public key: " + err.Error())
	}

	roleRepository := repositories.NewPostgresRoleRepository(db.DB)
	userRepository := repositories.NewPostgresUserRepository(db.DB)

	return &App{
		UserService: service.NewUserService(userRepository, roleRepository, keys),
		RoleService: service.NewRoleService(roleRepository),
		Keys:        keys,
	}
}
