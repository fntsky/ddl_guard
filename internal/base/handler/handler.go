package handler

import (
	"errors"
	"net/http"
	"strings"

	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/gin-gonic/gin"
)

// AppError 别名，保持兼容
type AppError = apperrors.AppError

// NewError 用于创建未预定义的错误
func NewError(code int, message string, err error) *AppError {
	return &AppError{HTTPStatus: code, Message: message, Err: err}
}

func BadRequest(message string, err error) *AppError {
	return NewError(http.StatusBadRequest, message, err)
}

func NotFound(message string, err error) *AppError {
	return NewError(http.StatusNotFound, message, err)
}

func Conflict(message string, err error) *AppError {
	return NewError(http.StatusConflict, message, err)
}

func Internal(message string, err error) *AppError {
	return NewError(http.StatusInternalServerError, message, err)
}

func BindAndCheck(ctx *gin.Context, data any) bool {
	if err := ctx.ShouldBindJSON(data); err != nil {
		HandleResponse(ctx, BadRequest("request format error", err), nil)
		return true
	}
	return false
}

func HandleResponse(ctx *gin.Context, err error, data any) {
	if err == nil {
		ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "Success", data))
		return
	}

	appErr := NormalizeError(err)
	message := strings.TrimSpace(appErr.Message)
	if message == "" {
		message = appErr.Error()
	}
	ctx.JSON(appErr.HTTPStatus, NewRespBodyData(appErr.HTTPStatus, message, data))
}

func NormalizeError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		// 确保 HTTPStatus 有默认值
		if appErr.HTTPStatus == 0 {
			appErr.HTTPStatus = http.StatusInternalServerError
		}
		if strings.TrimSpace(appErr.Message) == "" {
			appErr.Message = appErr.Error()
		}
		return appErr
	}

	// 未知错误返回 500
	return Internal("internal server error", err)
}
