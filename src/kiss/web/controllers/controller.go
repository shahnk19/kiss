package controllers

import (
	"github.com/gin-gonic/gin"
)

func Encode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": true})
	}
}
