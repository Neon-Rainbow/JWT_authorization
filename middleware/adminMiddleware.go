package middleware

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/controller"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		isAdmin, exist := c.Get("isAdmin")
		if !exist {
			controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Unauthorized access, don't exist isAdmin field in context")
			c.Abort()
			return
		}

		if isAdmin != true {
			controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Unauthorized access, admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}
