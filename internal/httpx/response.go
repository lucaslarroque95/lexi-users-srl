// Package httpx builds API Gateway v2 (HTTP API) responses in the same
// shapes the gin handlers in lexi-users-ms/src/routes already return.
package httpx

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// JSON marshals body and wraps it in a 200-family (or any status) response.
func JSON(status int, body any) (events.APIGatewayV2HTTPResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return Error(500, "could not encode response")
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(payload),
	}, nil
}

// Empty returns a body-less response (mirrors gin's c.Status, e.g. 204).
func Empty(status int) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{StatusCode: status}, nil
}

// Error mirrors gin.H{"error": msg} — the shape routes/*.go uses for every
// failure except Login, which uses "message" instead (see MessageError).
func Error(status int, msg string) (events.APIGatewayV2HTTPResponse, error) {
	return JSON(status, map[string]string{"error": msg})
}

// MessageError mirrors gin.H{"message": msg}, used by Login on failure.
func MessageError(status int, msg string) (events.APIGatewayV2HTTPResponse, error) {
	return JSON(status, map[string]string{"message": msg})
}
