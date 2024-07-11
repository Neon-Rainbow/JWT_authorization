package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    code.ResponseCode `json:"code"`
	Message string            `json:"message"`
	Error   error             `json:"error"`
	Data    interface{}       `json:"data"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	response := Response{
		Code:    code.Success,
		Message: code.Success.Message(),
		Data:    data,
	}
	c.JSON(200, response)
	return
}

func ResponseErrorWithCode(c *gin.Context, code code.ResponseCode) {
	response := Response{
		Code:    code,
		Message: code.Message(),
		Error:   nil,
		Data:    nil,
	}
	c.JSON(200, response)
	return
}

func ResponseErrorWithApiError(c *gin.Context, apiError *model.ApiError) {
	response := Response{
		Code:    apiError.Code,
		Message: apiError.Message,
		Error:   apiError.ErrorMessage,
		Data:    nil,
	}
	c.JSON(http.StatusOK, response)
}

func ResponseErrorWithMessage(c *gin.Context, code code.ResponseCode, message string) {
	var msg string
	if message == "" {
		msg = code.Message()
	} else {
		msg = message
	}
	response := Response{
		Code:    code,
		Message: msg,
		Error:   nil,
		Data:    nil,
	}
	c.JSON(200, response)
	return
}

// ResponseWithHttpStatus is a function that returns a response with a specific HTTP status code
func ResponseWithHttpStatus(c *gin.Context, httpStatus int, code code.ResponseCode, message string) {
	var msg string
	if message == "" {
		msg = code.Message()
	} else {
		msg = message
	}
	response := Response{
		Code:    code,
		Message: msg,
		Error:   nil,
		Data:    nil,
	}
	c.JSON(httpStatus, response)
	return
}
