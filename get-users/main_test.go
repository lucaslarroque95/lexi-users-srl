package main

import (
	"encoding/json"
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

func TestHandler_FiltersByRoleQueryParam(t *testing.T) {
	app, userRepo := newTestApp(t)
	userRepo.Users["1"] = models.User{ID: "1", Email: "admin@lexi.com", Roles: []models.Role{{Name: "admin"}}}
	userRepo.Users["2"] = models.User{ID: "2", Email: "user@lexi.com", Roles: []models.Role{{Name: "default"}}}

	req := events.APIGatewayV2HTTPRequest{
		QueryStringParameters: map[string]string{"role": "admin"},
		Headers:               map[string]string{"authorization": adminToken(t, app)},
	}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}

	var users []map[string]any
	if err := json.Unmarshal([]byte(resp.Body), &users); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(users) != 1 || users[0]["email"] != "admin@lexi.com" {
		t.Fatalf("expected only the admin user, got %v", users)
	}
}

func TestHandler_ForbiddenWithoutAdminRole(t *testing.T) {
	app, userRepo := newTestApp(t)
	userRepo.Users["1"] = models.User{ID: "1", Email: "user@lexi.com"}

	userToken, err := app.Keys.GenerateToken("1", "user@lexi.com", []string{"default"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := events.APIGatewayV2HTTPRequest{Headers: map[string]string{"authorization": userToken}}
	resp, err := Handler(app, t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}
