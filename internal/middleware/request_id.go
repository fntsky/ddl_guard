package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// RequestIDMiddleware uses client's X-Request-ID if present, otherwise generates a new UUID.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set(RequestIDKey, requestID)
		ctx.Header(RequestIDHeader, requestID)

		ctx.Next()
	}
}

func GetRequestID(ctx *gin.Context) string {
	return ctx.GetString(RequestIDKey)
}