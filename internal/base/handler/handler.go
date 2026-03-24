package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindAndCheck(ctx *gin.Context, data any) bool {
	if err := ctx.ShouldBindJSON(data); err != nil {
		return true
	}
	return false
}

func HanderResponse(ctx *gin.Context, err error, data any) {
	if err == nil {
		ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "Success", data))
	}
	ctx.JSON(http.StatusInternalServerError, NewRespBodyData(http.StatusInternalServerError, err.Error(), nil))
}
