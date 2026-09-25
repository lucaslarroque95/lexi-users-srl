package main

import (
	"encoding/json"
	"testing"

	"lexi/users/models"
	"lexi/users/service"
	"lexi/users/test/testutil"
	"lexi/users/utils"

	"lexi-users-srl/internal/bootstrap"

	"github.com/aws/aws-lambda-go/events"
)

func newTestApp(t *testing.T) (*bootstrap.App, *testutil.FakeUserRepository) {
	t.Helper()

	userRepo := testutil.NewFakeUserRepository()
	roleRepo := testutil.NewFakeRoleRepository()
	roleRepo.Create(models.Role{Name: "default"})

	keys := testutil.NewKeys(t)
	return &bootstrap.App{
		UserService: service.NewUserService(userRepo, roleRepo, keys),
		RoleService: service.NewRoleService(roleRepo),
		Keys:        keys,
	}, userRepo
}

func TestHandler_Success(t *testing.T) {
	app, userRepo := newTestApp(t)
	hashed, err := utils.HashPassword("secret")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	userRepo.Users["1"] = models.User{ID: "1", Email: "admin@lexi.com", Password: hashed, Roles: []models.Role{{Name: "admin"}}}

	req := events.APIGatewayV2HTTPRequest{Body: `{"email":"admin@lexi.com","password":"secret"}`}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}

	var body map[string]any
	if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	token, _ := body["token"].(string)
	if token == "" {
		t.Fatalf("expected a non-empty token in the response, got %v", body)
	}
}

func TestHandler_WrongPasswordReturns401(t *testing.T) {
	app, userRepo := newTestApp(t)
	hashed, _ := utils.HashPassword("secret")
	userRepo.Users["1"] = models.User{ID: "1", Email: "admin@lexi.com", Password: hashed}

	req := events.APIGatewayV2HTTPRequest{Body: `{"email":"admin@lexi.com","password":"wrong"}`}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
