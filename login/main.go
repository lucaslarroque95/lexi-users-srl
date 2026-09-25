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

// Handler mirrors routes.UserHandler.Login (POST /login, public).
func Handler(app *bootstrap.App, ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var payload schemas.UserLogin
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpx.Error(400, err.Error())
	}

	token, err := app.UserService.LogIn(models.User{Email: payload.Email, Password: payload.Password})
	if err != nil {
		return httpx.MessageError(401, err.Error())
	}

	return httpx.JSON(200, schemas.LoginResponse{Message: "User authenticated successfully", Token: token})
}

func main() {
	app := bootstrap.New()
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return Handler(app, ctx, req)
	})
}
