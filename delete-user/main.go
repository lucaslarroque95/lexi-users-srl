package main

import (
	"context"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.UserHandler.DeleteUser (DELETE /users/{id}, admin only).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	_, roles, err := httpx.Authenticate(app.Keys, req)
	if err != nil {
		return httpx.MessageError(401, "Not authorized")
	}
	if !httpx.RequireRole(roles, "admin") {
		return httpx.MessageError(403, "Forbidden")
	}

	id := req.PathParameters["id"]

	if err := app.UserService.DeleteUser(id); err != nil {
		return httpx.Error(500, "could not delete user")
	}

	return httpx.Empty(204)
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
