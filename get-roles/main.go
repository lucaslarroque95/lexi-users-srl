package main

import (
	"context"

	"lexi-users-srl/internal/lib/schemas"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.RoleHandler.GetRoles (GET /roles, admin only).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	_, roles, err := httpx.Authenticate(app.Keys, req)
	if err != nil {
		return httpx.MessageError(401, "Not authorized")
	}
	if !httpx.RequireRole(roles, "admin") {
		return httpx.MessageError(403, "Forbidden")
	}

	all, err := app.RoleService.ListRoles()
	if err != nil {
		return httpx.Error(500, "could not fetch roles")
	}

	result := make([]schemas.RoleRead, len(all))
	for i, r := range all {
		result[i] = schemas.RoleRead{ID: r.ID, Name: r.Name}
	}
	return httpx.JSON(200, result)
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
