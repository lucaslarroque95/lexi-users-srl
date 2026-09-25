package schemas

// ErrorResponse is the shape of most error responses returned by this API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse is returned by Login on failure (uses "message", not "error").
type MessageResponse struct {
	Message string `json:"message"`
}

// LoginResponse is returned by a successful Login.
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
