package main

import (
	"context"
	"encoding/json"
	"strconv"

	"lexi-users-srl/internal/lib/schemas"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.RoleHandler.UpdateRole (PUT /roles/{id}, admin only).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	_, roles, err := httpx.Authenticate(app.Keys, req)
	if err != nil {
		return httpx.MessageError(401, "Not authorized")
	}
	if !httpx.RequireRole(roles, "admin") {
		return httpx.MessageError(403, "Forbidden")
	}

	id, err := strconv.ParseInt(req.PathParameters["id"], 10, 64)
	if err != nil {
		return httpx.Error(400, "invalid role id")
	}

	var payload schemas.RoleUpdate
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpx.Error(400, err.Error())
	}

	existing, err := app.RoleService.GetRole(id)
	if err != nil {
		return httpx.Error(404, "role not found")
	}

	if payload.Name != nil {
		existing.Name = *payload.Name
	}

	updated, err := app.RoleService.UpdateRole(id, existing)
	if err != nil {
		return httpx.Error(500, "could not update role")
	}

	return httpx.JSON(200, schemas.RoleRead{ID: updated.ID, Name: updated.Name})
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
