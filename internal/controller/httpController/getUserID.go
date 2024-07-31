package httpController

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (userID string, exists bool) {
	isAdmin, _ := c.Get("isAdmin")
	isAdmin = isAdmin.(bool)
	if !isAdmin.(bool) {
		_any_type_userID, exist := c.Get("userID")
		if !exist {
			return "", false
		}
		userID = fmt.Sprintf("%v", _any_type_userID) // _any_type_userID is interface{} type, so we need to convert it to string type to get the value of
	} else {
		userID = c.Query("user_id")
		if userID == "" {
			return "", false
		}
	}
	return userID, true
}
