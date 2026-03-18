package handler

import "github.com/gin-gonic/gin"

func BindAndCheck(ctx *gin.Context, data any) bool {
	if err := ctx.ShouldBindJSON(data); err != nil {
		return true
	}
	return false
}
