package middleware

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		isAdmin, exist := c.Get("isAdmin")
		if !exist {
			controller.ResponseWithHttpStatus(c, http.StatusUnauthorized, code.RequestUnauthorized, "Unauthorized access, don't exist isAdmin field in context")
			//controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Unauthorized access, don't exist isAdmin field in context")
			c.Abort()
			return
		}

		if isAdmin != true {
			controller.ResponseWithHttpStatus(c, http.StatusUnauthorized, code.RequestUnauthorized, "Unauthorized access, admin only")
			controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Unauthorized access, admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}
