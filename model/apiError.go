package model

import (
	"JWT_authorization/code"
	"fmt"
)

type ApiError struct {
	Code         code.ResponseCode `json:"code"`
	Message      string            `json:"message"`
	ErrorMessage error             `json:"error"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("ApiError Detail | Code: %v | Message: %v | ErrorMessage: %v", e.Code, e.Message, e.ErrorMessage)
}
