package main

import (
	"encoding/json"
	"testing"

	"lexi/users/models"
	"lexi/users/service"
	"lexi/users/test/testutil"

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

func TestHandler_CreatesUser(t *testing.T) {
	app, userRepo := newTestApp(t)

	req := events.APIGatewayV2HTTPRequest{Body: `{"email":"new@lexi.com","password":"secret"}`}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, resp.Body)
	}
	if len(userRepo.Users) != 1 {
		t.Fatalf("expected 1 user to be persisted, got %d", len(userRepo.Users))
	}

	var body map[string]any
	if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["email"] != "new@lexi.com" {
		t.Fatalf("expected email in response, got %v", body)
	}
	if _, exposesPassword := body["password"]; exposesPassword {
		t.Fatalf("expected the password not to be exposed in the response")
	}
}

func TestHandler_InvalidBodyReturns400(t *testing.T) {
	app, _ := newTestApp(t)

	req := events.APIGatewayV2HTTPRequest{Body: `not-json`}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
