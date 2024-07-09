package controller

import (
	"github.com/gin-gonic/gin"
)

func GetAdminInfo(c *gin.Context) {
	ResponseSuccess(c, gin.H{
		"message": "get admin info success",
	})
}
