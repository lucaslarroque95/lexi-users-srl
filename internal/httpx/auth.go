package httpx

import (
	"errors"
	"slices"

	"lexi/users/utils"

	"github.com/aws/aws-lambda-go/events"
)

// Authenticate mirrors middlewares.Authenticate: reads the raw JWT from the
// Authorization header and verifies it against keys. API Gateway v2 lower-cases
// header names, so it only needs to check "authorization".
func Authenticate(keys utils.Keys, req events.APIGatewayV2HTTPRequest) (userID string, roles []string, err error) {
	token := req.Headers["authorization"]
	if token == "" {
		return "", nil, errors.New("not authorized")
	}

	return keys.VerifyToken(token)
}

// RequireRole mirrors middlewares.RequireRole.
func RequireRole(roles []string, role string) bool {
	return slices.Contains(roles, role)
}
