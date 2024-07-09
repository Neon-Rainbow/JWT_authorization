package model

import "JWT_authorization/code"

type ApiError struct {
	Code    code.ResponseCode `json:"code"`
	Message string            `json:"message"`
	Error   error             `json:"error"`
}
