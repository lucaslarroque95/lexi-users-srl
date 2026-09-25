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

// Handler mirrors routes.UserHandler.UpdateOwnUser (PUT /user, authenticated;
// unlike PUT /users/{id}, this one cannot reassign roles).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID, _, err := httpx.Authenticate(app.Keys, req)
	if err != nil {
		return httpx.MessageError(401, "Not authorized")
	}

	var payload schemas.UserUpdate
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpx.Error(400, err.Error())
	}

	existing, err := app.UserService.GetUser(userID)
	if err != nil {
		return httpx.Error(404, "user not found")
	}

	if payload.Email != nil {
		existing.Email = *payload.Email
	}
	if payload.Password != nil {
		existing.Password = *payload.Password
	}

	updated, err := app.UserService.UpdateUser(userID, existing)
	if err != nil {
		return httpx.Error(500, "could not update user")
	}

	return httpx.JSON(200, schemas.UserRead{ID: updated.ID, Email: updated.Email})
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
