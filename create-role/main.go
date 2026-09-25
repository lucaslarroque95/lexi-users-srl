package main

import (
	"context"
	"encoding/json"

	"lexi/users/models"
	"lexi/users/schemas"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.RoleHandler.CreateRole (POST /roles, admin only).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	_, roles, err := httpx.Authenticate(app.Keys, req)
	if err != nil {
		return httpx.MessageError(401, "Not authorized")
	}
	if !httpx.RequireRole(roles, "admin") {
		return httpx.MessageError(403, "Forbidden")
	}

	var payload schemas.RoleCreate
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpx.Error(400, err.Error())
	}

	created, err := app.RoleService.CreateRole(models.Role{Name: payload.Name})
	if err != nil {
		return httpx.Error(500, "could not create role")
	}

	return httpx.JSON(201, schemas.RoleRead{ID: created.ID, Name: created.Name})
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
