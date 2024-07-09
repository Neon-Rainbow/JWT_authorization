package controller

import (
	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	ResponseSuccess(c, gin.H{
		"message": "Get user info successfully",
	})
}
