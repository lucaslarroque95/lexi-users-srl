package main

import (
	"context"

	"lexi/users/models"
	"lexi/users/schemas"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.UserHandler.GetUsers (GET /users, admin only,
// optionally filtered by the "role" query param).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	_, roles, err := httpx.Authenticate(app.Keys, req)
	if err != nil {
		return httpx.MessageError(401, "Not authorized")
	}
	if !httpx.RequireRole(roles, "admin") {
		return httpx.MessageError(403, "Forbidden")
	}

	role := req.QueryStringParameters["role"]

	var users []models.User
	if role != "" {
		users, err = app.UserService.ListUsersByRole(role)
	} else {
		users, err = app.UserService.ListUsers()
	}
	if err != nil {
		return httpx.Error(500, "could not fetch users")
	}

	result := make([]schemas.UserRead, len(users))
	for i, u := range users {
		result[i] = schemas.UserRead{ID: u.ID, Email: u.Email}
	}
	return httpx.JSON(200, result)
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
