package controller

import (
	"JWT_authorization/code"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) string {
	isAdmin, _ := c.Get("isAdmin")
	var userID string
	if !isAdmin.(bool) {
		_any_type_userID, exist := c.Get("userID")
		if !exist {
			ResponseErrorWithCode(c, code.ServerBusy)
			return ""
		}
		userID = fmt.Sprintf("%v", _any_type_userID) // _any_type_userID is interface{} type, so we need to convert it to string type to get the value of
	} else {
		userID = c.Query("user_id")
		if userID == "" {
			ResponseErrorWithCode(c, code.ServerBusy)
			return ""
		}
	}
	return userID
}
