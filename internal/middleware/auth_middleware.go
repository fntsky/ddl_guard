package middleware

import (
	"errors"
	"strings"

	baseauth "github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/base/handler"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	pkgjwt "github.com/fntsky/ddl_guard/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const GinAuthUserKey = "auth_user"

func AuthMiddleware(tokenService *baseauth.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			handler.HandleResponse(ctx, apperrors.ErrTokenMissing, nil)
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			handler.HandleResponse(ctx, apperrors.ErrTokenMalformed, nil)
			ctx.Abort()
			return
		}

		claims, err := tokenService.ParseAccessToken(strings.TrimSpace(parts[1]))
		if err != nil {
			// 将 pkgjwt 错误转换为 apperrors
			var appErr *apperrors.AppError
			if errors.As(err, &appErr) {
				handler.HandleResponse(ctx, appErr, nil)
			} else if errors.Is(err, pkgjwt.ErrTokenExpired) {
				handler.HandleResponse(ctx, apperrors.ErrTokenExpired, nil)
			} else if errors.Is(err, pkgjwt.ErrInvalidToken) || errors.Is(err, pkgjwt.ErrInvalidSignature) || errors.Is(err, pkgjwt.ErrInvalidClaims) {
				handler.HandleResponse(ctx, apperrors.ErrInvalidToken, nil)
			} else {
				handler.HandleResponse(ctx, apperrors.ErrInvalidToken, nil)
			}
			ctx.Abort()
			return
		}
		ctx.Set(GinAuthUserKey, claims)
		ctx.Next()
	}
}

func GetUserFromGin(ctx *gin.Context) (*pkgjwt.Claims, bool) {
	v, ok := ctx.Get(GinAuthUserKey)
	if !ok {
		return nil, false
	}
	claims, ok := v.(*pkgjwt.Claims)
	return claims, ok
}
