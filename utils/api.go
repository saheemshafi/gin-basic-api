package utils

import (
	"github.com/gin-gonic/gin"
)

func WriteResponse(ctx *gin.Context, status int, message string, data ...any) {
	var response any

	if len(data) == 1 {
		response = data[0]
	}

	ctx.JSON(status, gin.H{
		"status":  status,
		"message": message,
		"data":    response,
	})
}
