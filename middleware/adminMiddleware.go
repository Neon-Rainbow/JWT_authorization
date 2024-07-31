package middleware

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/controller/httpController"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func AdminMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		isAdmin, exist := c.Get("isAdmin")
		if !exist {
			httpController.ResponseWithHttpStatus(
				c,
				http.StatusUnauthorized,
				code.RequestUnauthorized,
				"Unauthorized access, don't exist isAdmin field in context",
			)
			zap.L().Error("AdminMiddleware: Unauthorized access, don't exist isAdmin field in context",
				zap.Int("code", code.RequestUnauthorized.Int()),
			)

			//httpController.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Unauthorized access, don't exist isAdmin field in context")
			c.Abort()
			return
		}

		if isAdmin != true {
			httpController.ResponseWithHttpStatus(
				c,
				http.StatusUnauthorized,
				code.RequestUnauthorized,
				"Unauthorized access, admin only",
			)
			zap.L().Error("AdminMiddleware: Unauthorized access, admin only",
				zap.Int("code", code.RequestUnauthorized.Int()),
			)
			// httpController.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Unauthorized access, admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}
