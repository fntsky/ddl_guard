package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ErrorCode 机器可读的错误码，供客户端程序化处理
type ErrorCode string

// 错误码常量
const (
	// 成功
	CodeSuccess ErrorCode = "SUCCESS"

	// DDL 领域错误
	CodeDraftNotFound      ErrorCode = "DRAFT_NOT_FOUND"
	CodeInvalidDraftStatus ErrorCode = "INVALID_DRAFT_STATUS"
	CodeDraftStateConflict ErrorCode = "DRAFT_STATE_CONFLICT"
	CodePictureDataMissing ErrorCode = "PICTURE_DATA_MISSING"
	CodePictureDataInvalid ErrorCode = "PICTURE_DATA_INVALID"
	CodeDeadlineInPast     ErrorCode = "DEADLINE_IN_PAST"
	CodeDDLNotFound        ErrorCode = "DDL_NOT_FOUND"
	CodeDDLNotActive       ErrorCode = "DDL_NOT_ACTIVE"
	CodeDDLNotOwned        ErrorCode = "DDL_NOT_OWNED"
	CodeAIProviderDisabled ErrorCode = "AI_PROVIDER_DISABLED"

	// User 领域错误
	CodeUserNotFound                ErrorCode = "USER_NOT_FOUND"
	CodeEmailAlreadyExists          ErrorCode = "EMAIL_ALREADY_EXISTS"
	CodeInvalidCredentials          ErrorCode = "INVALID_CREDENTIALS"
	CodeInvalidVerificationCode     ErrorCode = "INVALID_VERIFICATION_CODE"
	CodeVerificationUnavailable     ErrorCode = "VERIFICATION_UNAVAILABLE"
	CodeEmailOTPDisabled            ErrorCode = "EMAIL_OTP_DISABLED"
	CodeUnsupportedVerificationType ErrorCode = "UNSUPPORTED_VERIFICATION_TYPE"

	// Auth 领域错误
	CodeTokenMissing        ErrorCode = "TOKEN_MISSING"
	CodeTokenMalformed      ErrorCode = "TOKEN_MALFORMED"
	CodeInvalidToken        ErrorCode = "INVALID_TOKEN"
	CodeTokenExpired        ErrorCode = "TOKEN_EXPIRED"
	CodeInvalidRefreshToken ErrorCode = "INVALID_REFRESH_TOKEN"
	CodeRefreshTokenExpired ErrorCode = "REFRESH_TOKEN_EXPIRED"
	CodeRefreshTokenRevoked ErrorCode = "REFRESH_TOKEN_REVOKED"
	CodeSessionNotFound     ErrorCode = "SESSION_NOT_FOUND"
	CodeTokenConfigInvalid  ErrorCode = "TOKEN_CONFIG_INVALID"

	// 基础设施错误
	CodeDatabaseError ErrorCode = "DATABASE_ERROR"
	CodeRedisError    ErrorCode = "REDIS_ERROR"
	CodeInternalError ErrorCode = "INTERNAL_ERROR"

	// 通用 HTTP 错误
	CodeBadRequest   ErrorCode = "BAD_REQUEST"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeConflict     ErrorCode = "CONFLICT"

	// AI 提供者错误
	CodeAIRequestFailed   ErrorCode = "AI_REQUEST_FAILED"
	CodeAIResponseInvalid ErrorCode = "AI_RESPONSE_INVALID"
)

// AppError 统一错误类型
type AppError struct {
	HTTPStatus int       `json:"-"`       // HTTP 状态码
	Code       ErrorCode `json:"code"`    // 机器可读错误码
	Message    string    `json:"message"` // 用户友好消息
	Err        error     `json:"-"`       // 原始错误（不序列化）

	// 可选上下文字段
	RequestID string `json:"request_id,omitempty"` // 请求追踪 ID
	UserID    string `json:"user_id,omitempty"`    // 用户 ID
	Operation string `json:"-"`                    // 操作名称，用于日志

	// 堆栈跟踪（创建时捕获，不序列化到 JSON）
	stack []uintptr `json:"-"`
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(string(e.Code))
	b.WriteString(": ")
	b.WriteString(e.Message)
	if e.Operation != "" {
		b.WriteString(" [op=")
		b.WriteString(e.Operation)
		b.WriteString("]")
	}
	if e.Err != nil {
		b.WriteString(": ")
		b.WriteString(e.Err.Error())
	}
	return b.String()
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *AppError) StackTrace() []runtime.Frame {
	if e == nil || len(e.stack) == 0 {
		return nil
	}
	frames := runtime.CallersFrames(e.stack)
	var result []runtime.Frame
	for {
		frame, more := frames.Next()
		result = append(result, frame)
		if !more {
			break
		}
	}
	return result
}

func (e *AppError) HasStack() bool {
	return e != nil && len(e.stack) > 0
}

func (e *AppError) WithRequestID(requestID string) *AppError {
	if e == nil {
		return nil
	}
	e.RequestID = requestID
	return e
}

func (e *AppError) WithUserID(userID string) *AppError {
	if e == nil {
		return nil
	}
	e.UserID = userID
	return e
}

func (e *AppError) WithOperation(op string) *AppError {
	if e == nil {
		return nil
	}
	e.Operation = op
	return e
}

func (e *AppError) Wrap(err error) *AppError {
	if e == nil || err == nil {
		return e
	}
	return &AppError{
		HTTPStatus: e.HTTPStatus,
		Code:       e.Code,
		Message:    e.Message,
		Err:        err,
		RequestID:  e.RequestID,
		UserID:     e.UserID,
		Operation:  e.Operation,
		stack:      e.stack,
	}
}

func New(httpStatus int, code ErrorCode, message string) *AppError {
	const depth = 3
	var pcs [32]uintptr
	n := runtime.Callers(depth, pcs[:])
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		stack:      pcs[:n],
	}
}

func WrapError(err error, httpStatus int, code ErrorCode, message string) *AppError {
	appErr := New(httpStatus, code, message)
	appErr.Err = err
	return appErr
}

// 预定义错误
var (
	// DDL 相关
	ErrInvalidDraftStatus = New(http.StatusBadRequest, CodeInvalidDraftStatus, "invalid draft status")
	ErrDraftNotFound      = New(http.StatusNotFound, CodeDraftNotFound, "draft not found")
	ErrDraftStateConflict = New(http.StatusConflict, CodeDraftStateConflict, "draft state conflict")
	ErrPictureDataMissing = New(http.StatusBadRequest, CodePictureDataMissing, "picture base64 data is required")
	ErrPictureDataInvalid = New(http.StatusBadRequest, CodePictureDataInvalid, "invalid picture base64 data")
	ErrAIProviderDisabled = New(http.StatusInternalServerError, CodeAIProviderDisabled, "ai provider is not configured")
	ErrUserNotFound       = New(http.StatusNotFound, CodeUserNotFound, "user not found")
	ErrDraftNotOwned      = New(http.StatusForbidden, CodeDDLNotOwned, "draft not owned by user")
	ErrDDLNotFound        = New(http.StatusNotFound, CodeDDLNotFound, "ddl not found")
	ErrDDLNotActive       = New(http.StatusBadRequest, CodeDDLNotActive, "ddl is not active")
	ErrDeadlineInPast     = New(http.StatusBadRequest, CodeDeadlineInPast, "deadline cannot be earlier than current time")

	// User 相关
	ErrEmailAlreadyExists          = New(http.StatusConflict, CodeEmailAlreadyExists, "email already exists")
	ErrInvalidVerificationCode     = New(http.StatusBadRequest, CodeInvalidVerificationCode, "invalid verification code")
	ErrVerificationUnavailable     = New(http.StatusInternalServerError, CodeVerificationUnavailable, "verification service is not available")
	ErrEmailOTPDisabled            = New(http.StatusServiceUnavailable, CodeEmailOTPDisabled, "email otp is disabled")
	ErrInvalidCredentials          = New(http.StatusUnauthorized, CodeInvalidCredentials, "invalid email or password")
	ErrUnsupportedVerificationType = New(http.StatusBadRequest, CodeUnsupportedVerificationType, "unsupported verification type")

	// Auth 相关
	ErrInvalidRefreshToken = New(http.StatusUnauthorized, CodeInvalidRefreshToken, "invalid refresh token")
	ErrRefreshTokenExpired = New(http.StatusUnauthorized, CodeRefreshTokenExpired, "refresh token expired")
	ErrRefreshTokenRevoked = New(http.StatusUnauthorized, CodeRefreshTokenRevoked, "refresh token revoked")
	ErrSessionNotFound     = New(http.StatusUnauthorized, CodeSessionNotFound, "refresh session not found")
	ErrTokenMissing        = New(http.StatusUnauthorized, CodeTokenMissing, "authorization token is required")
	ErrTokenMalformed      = New(http.StatusUnauthorized, CodeTokenMalformed, "authorization header must be Bearer token")
	ErrInvalidToken        = New(http.StatusUnauthorized, CodeInvalidToken, "invalid token")
	ErrTokenExpired        = New(http.StatusUnauthorized, CodeTokenExpired, "token expired")
	ErrTokenConfigInvalid  = New(http.StatusInternalServerError, CodeTokenConfigInvalid, "token config is invalid")

	// OTP 相关
	ErrCodeStoreNotConfigured = New(http.StatusInternalServerError, CodeRedisError, "code store not configured")

	// Exam 相关
	ErrExamNotFound    = New(http.StatusNotFound, CodeNotFound, "exam not found")
	ErrExamNotOwned    = New(http.StatusForbidden, CodeForbidden, "exam not owned by user")
	ErrExamTimeInvalid = New(http.StatusBadRequest, CodeBadRequest, "end time must be after start time")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func AIRequestFailed(provider string, statusCode int, body string, err error) *AppError {
	return WrapError(
		err,
		http.StatusBadGateway,
		CodeAIRequestFailed,
		fmt.Sprintf("AI provider '%s' request failed", provider),
	).WithOperation("ai.request")
}

func AIResponseInvalid(provider string, reason string, err error) *AppError {
	return WrapError(
		err,
		http.StatusBadGateway,
		CodeAIResponseInvalid,
		fmt.Sprintf("AI provider '%s' returned invalid response: %s", provider, reason),
	).WithOperation("ai.response")
}

func DatabaseError(operation string, err error) *AppError {
	return WrapError(
		err,
		http.StatusInternalServerError,
		CodeDatabaseError,
		"database operation failed",
	).WithOperation(operation)
}

func RedisError(operation string, err error) *AppError {
	return WrapError(
		err,
		http.StatusInternalServerError,
		CodeRedisError,
		"redis operation failed",
	).WithOperation(operation)
}

func ValidationError(field string, reason string) *AppError {
	return New(http.StatusBadRequest, CodeBadRequest,
		fmt.Sprintf("validation failed for '%s': %s", field, reason),
	).WithOperation("validation")
}

func NotFoundError(resource string, identifier string) *AppError {
	return New(http.StatusNotFound, CodeNotFound,
		fmt.Sprintf("%s not found", resource),
	).WithOperation(fmt.Sprintf("%s.lookup:%s", resource, identifier))
}

func UnauthorizedError(reason string) *AppError {
	return New(http.StatusUnauthorized, CodeUnauthorized, reason).
		WithOperation("auth.check")
}

func ForbiddenError(resource string) *AppError {
	return New(http.StatusForbidden, CodeForbidden,
		fmt.Sprintf("access denied to %s", resource),
	).WithOperation("auth.authorize")
}
