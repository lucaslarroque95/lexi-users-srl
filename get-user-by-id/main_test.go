package main

import (
	"testing"

	"lexi-users-srl/internal/lib/models"
	"lexi-users-srl/internal/lib/service"
	"lexi-users-srl/internal/lib/test/testutil"

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

func adminToken(t *testing.T, app *bootstrap.App) string {
	t.Helper()
	token, err := app.Keys.GenerateToken("admin-id", "admin@lexi.com", []string{"admin"})
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}
	return token
}

func TestHandler_NotFoundReturns404(t *testing.T) {
	app, _ := newTestApp(t)

	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "999"},
		Headers:        map[string]string{"authorization": adminToken(t, app)},
	}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Fatalf("expected 404, got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestHandler_NoTokenReturns401(t *testing.T) {
	app, _ := newTestApp(t)

	req := events.APIGatewayV2HTTPRequest{PathParameters: map[string]string{"id": "1"}}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestHandler_NonAdminReturns403(t *testing.T) {
	app, userRepo := newTestApp(t)
	userRepo.Users["1"] = models.User{ID: "1", Email: "user@lexi.com"}

	userToken, err := app.Keys.GenerateToken("1", "user@lexi.com", []string{"default"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "1"},
		Headers:        map[string]string{"authorization": userToken},
	}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}
