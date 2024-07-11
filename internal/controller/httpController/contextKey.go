package httpController

type contextKey string

const (
	apiErrorKey         contextKey = "api_error"
	loginResponseKey    contextKey = "login_response"
	registerResponseKey contextKey = "register_response"
)
