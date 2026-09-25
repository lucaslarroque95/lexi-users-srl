package main

import (
	"context"
	"encoding/json"

	"lexi-users-srl/internal/lib/models"
	"lexi-users-srl/internal/lib/schemas"

	"lexi-users-srl/internal/bootstrap"
	"lexi-users-srl/internal/httpx"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler mirrors routes.UserHandler.SignUp (POST /signup, public).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var payload schemas.UserCreate
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpx.Error(400, err.Error())
	}

	created, err := app.UserService.SignUp(models.User{Email: payload.Email, Password: payload.Password})
	if err != nil {
		return httpx.Error(500, "could not create user")
	}

	return httpx.JSON(201, schemas.UserRead{ID: created.ID, Email: created.Email})
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
