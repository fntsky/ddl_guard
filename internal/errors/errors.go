package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError 统一错误类型
type AppError struct {
	HTTPStatus int    // HTTP 状态码
	Message    string // 用户友好消息
	Err        error  // 原始错误（可选）
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// New 创建新错误
func New(httpStatus int, message string) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Message:    message,
	}
}

// Wrap 包装原始错误
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		HTTPStatus: e.HTTPStatus,
		Message:    e.Message,
		Err:        err,
	}
}

// 预定义错误
var (
	// DDL 相关
	ErrInvalidDraftStatus = New(http.StatusBadRequest, "invalid draft status")
	ErrDraftNotFound      = New(http.StatusNotFound, "draft not found")
	ErrDraftStateConflict = New(http.StatusConflict, "draft state conflict")
	ErrPictureDataMissing = New(http.StatusBadRequest, "picture base64 data is required")
	ErrPictureDataInvalid = New(http.StatusBadRequest, "invalid picture base64 data")
	ErrAIProviderDisabled = New(http.StatusInternalServerError, "ai provider is not configured")
	ErrUserNotFound       = New(http.StatusNotFound, "user not found")
	ErrDraftNotOwned      = New(http.StatusForbidden, "draft not owned by user")
	ErrDDLNotFound        = New(http.StatusNotFound, "ddl not found")
	ErrDDLNotActive       = New(http.StatusBadRequest, "ddl is not active")
	ErrDeadlineInPast     = New(http.StatusBadRequest, "deadline cannot be earlier than current time")

	// User 相关
	ErrEmailAlreadyExists      = New(http.StatusConflict, "email already exists")
	ErrInvalidVerificationCode = New(http.StatusBadRequest, "invalid verification code")
	ErrVerificationUnavailable = New(http.StatusInternalServerError, "verification service is not available")
	ErrEmailOTPDisabled        = New(http.StatusServiceUnavailable, "email otp is disabled")
	ErrInvalidCredentials      = New(http.StatusUnauthorized, "invalid email or password")

	// Auth 相关
	ErrInvalidRefreshToken = New(http.StatusUnauthorized, "invalid refresh token")
	ErrRefreshTokenExpired = New(http.StatusUnauthorized, "refresh token expired")
	ErrRefreshTokenRevoked = New(http.StatusUnauthorized, "refresh token revoked")
	ErrSessionNotFound     = New(http.StatusUnauthorized, "refresh session not found")
	ErrTokenMissing        = New(http.StatusUnauthorized, "authorization token is required")
	ErrTokenMalformed      = New(http.StatusUnauthorized, "authorization header must be Bearer token")
	ErrInvalidToken        = New(http.StatusUnauthorized, "invalid token")
	ErrTokenExpired        = New(http.StatusUnauthorized, "token expired")
	ErrTokenConfigInvalid  = New(http.StatusInternalServerError, "token config is invalid")

	// OTP 相关
	ErrCodeStoreNotConfigured      = New(http.StatusInternalServerError, "code store not configured")
	ErrUnsupportedVerificationType = New(http.StatusBadRequest, "unsupported verification type")

	// Exam 相关
	ErrExamNotFound     = New(http.StatusNotFound, "exam not found")
	ErrExamNotOwned     = New(http.StatusForbidden, "exam not owned by user")
	ErrExamTimeInvalid  = New(http.StatusBadRequest, "end time must be after start time")
)

// Is 提供错误比较支持，兼容标准库 errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 提供错误类型转换支持，兼容标准库 errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}
