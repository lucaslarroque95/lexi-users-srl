package main

import (
	"context"
	"strconv"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.RoleHandler.DeleteRole (DELETE /roles/{id}, admin only).
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

	if err := app.RoleService.DeleteRole(id); err != nil {
		return httpx.Error(500, "could not delete role")
	}

	return httpx.Empty(204)
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
