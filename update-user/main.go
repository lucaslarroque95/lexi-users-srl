package main

import (
	"context"
	"encoding/json"

	"lexi/users/schemas"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.UserHandler.UpdateUser (PUT /users/{id}, admin only;
// unlike PUT /user, this one can also reassign roles).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	_, roles, err := httpx.Authenticate(app.Keys, req)
	if err != nil {
		return httpx.MessageError(401, "Not authorized")
	}
	if !httpx.RequireRole(roles, "admin") {
		return httpx.MessageError(403, "Forbidden")
	}

	id := req.PathParameters["id"]

	var payload schemas.UserUpdate
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpx.Error(400, err.Error())
	}

	existing, err := app.UserService.GetUser(id)
	if err != nil {
		return httpx.Error(404, "user not found")
	}

	if payload.Email != nil {
		existing.Email = *payload.Email
	}
	if payload.Password != nil {
		existing.Password = *payload.Password
	}
	if payload.Roles != nil {
		resolved, err := app.UserService.ResolveRoles(*payload.Roles)
		if err != nil {
			return httpx.Error(400, "invalid role")
		}
		existing.Roles = resolved
	}

	updated, err := app.UserService.UpdateUser(id, existing)
	if err != nil {
		return httpx.Error(500, "could not update user")
	}

	roleNames := make([]string, len(updated.Roles))
	for i, r := range updated.Roles {
		roleNames[i] = r.Name
	}
	return httpx.JSON(200, schemas.UserReadByID{ID: updated.ID, Email: updated.Email, Role: roleNames})
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
