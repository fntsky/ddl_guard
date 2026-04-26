package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/gin-gonic/gin"
)

type AppError = apperrors.AppError

func NewError(code int, message string, err error) *AppError {
	return apperrors.WrapError(err, code, apperrors.CodeInternalError, message)
}

func BadRequest(message string, err error) *AppError {
	return apperrors.WrapError(err, http.StatusBadRequest, apperrors.CodeBadRequest, message)
}

func Unauthorized(message string, err error) *AppError {
	return apperrors.WrapError(err, http.StatusUnauthorized, apperrors.CodeUnauthorized, message)
}

func Forbidden(message string, err error) *AppError {
	return apperrors.WrapError(err, http.StatusForbidden, apperrors.CodeForbidden, message)
}

func NotFound(message string, err error) *AppError {
	return apperrors.WrapError(err, http.StatusNotFound, apperrors.CodeNotFound, message)
}

func Conflict(message string, err error) *AppError {
	return apperrors.WrapError(err, http.StatusConflict, apperrors.CodeConflict, message)
}

func Internal(message string, err error) *AppError {
	return apperrors.WrapError(err, http.StatusInternalServerError, apperrors.CodeInternalError, message)
}

func BindAndCheck(ctx *gin.Context, data any) bool {
	if err := ctx.ShouldBindJSON(data); err != nil {
		HandleResponse(ctx, BadRequest("request format error", err), nil)
		return true
	}
	return false
}

func HandleResponse(ctx *gin.Context, err error, data any) {
	requestID := ctx.GetString("request_id")

	if err == nil {
		ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "Success", data))
		return
	}

	appErr := NormalizeError(err)
	appErr.RequestID = requestID

	// auth middleware stores claims under "auth_user" key
	if v, ok := ctx.Get("auth_user"); ok {
		if claims, ok := v.(interface{ GetUserUUID() string }); ok {
			appErr.UserID = claims.GetUserUUID()
		}
	}

	logError(ctx, appErr)
	ctx.JSON(appErr.HTTPStatus, NewErrorResp(appErr, requestID))
}

func NormalizeError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus == 0 {
			appErr.HTTPStatus = http.StatusInternalServerError
		}
		if strings.TrimSpace(appErr.Message) == "" {
			appErr.Message = appErr.Error()
		}
		return appErr
	}

	return Internal("internal server error", err)
}

func logError(ctx *gin.Context, err *AppError) {
	logger := slog.With(
		"code", err.Code,
		"status", err.HTTPStatus,
		"operation", err.Operation,
		"request_id", err.RequestID,
		"path", ctx.Request.URL.Path,
		"method", ctx.Request.Method,
	)

	if err.UserID != "" {
		logger = logger.With("user_id", err.UserID)
	}

	switch {
	case err.HTTPStatus >= 500:
		if err.HasStack() {
			logger = logger.With("stack", formatStackTrace(err.StackTrace()))
		}
		if err.Err != nil {
			logger.Error(err.Message, "error", err.Err)
		} else {
			logger.Error(err.Message)
		}
	case err.HTTPStatus >= 400:
		logger.Info(err.Message)
	default:
		logger.Debug(err.Message)
	}
}

func formatStackTrace(frames []runtime.Frame) []string {
	var result []string
	for _, f := range frames {
		result = append(result, f.Function+" at "+f.File+":"+strconv.Itoa(f.Line))
	}
	return result
}