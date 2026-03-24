package handler

import (
	"errors"
	"net/http"
	"strings"

	serviceddl "github.com/fntsky/ddl_guard/internal/service/ddl"
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
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
	ctx.JSON(appErr.Code, NewRespBodyData(appErr.Code, message, data))
}

func NormalizeError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Code == 0 {
			appErr.Code = http.StatusInternalServerError
		}
		if strings.TrimSpace(appErr.Message) == "" {
			appErr.Message = appErr.Error()
		}
		return appErr
	}

	switch {
	case errors.Is(err, serviceddl.ErrInvalidDraftStatus):
		return BadRequest("invalid draft status", err)
	case errors.Is(err, serviceddl.ErrPictureDataMissing):
		return BadRequest("picture raw data is required", err)
	case errors.Is(err, serviceddl.ErrDraftNotFound):
		return NotFound("draft not found", err)
	case errors.Is(err, serviceddl.ErrDraftStateConflict):
		return Conflict("draft state conflict", err)
	case errors.Is(err, serviceddl.ErrAIProviderDisabled):
		return Internal("ai provider is not configured", err)
	default:
		return Internal("internal server error", err)
	}
}
