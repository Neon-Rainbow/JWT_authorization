package middleware

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/controller"
	"JWT_authorization/util/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWTMiddleware is a middleware that checks for a valid JWT token in the request header
func JWTMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tkn := parts[1]
		myClaims, err := jwt.ParseToken(tkn)
		if err != nil {
			controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		if myClaims.TokenType != "access_token" {
			controller.ResponseErrorWithMessage(c, code.RequestUnauthorized, "Invalid token type, must be access_token")
			c.Abort()
			return
		}

		c.Set("userID", myClaims.UserID)
		c.Set("username", myClaims.Username)
		c.Set("isAdmin", myClaims.IsAdmin)
		c.Next()
	}
}