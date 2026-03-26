package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	baseauth "github.com/fntsky/ddl_guard/internal/base/auth"
	pkgjwt "github.com/fntsky/ddl_guard/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const GinAuthUserKey = "auth_user"

var (
	ErrTokenMissing   = errors.New("authorization token is required")
	ErrTokenMalformed = errors.New("authorization header must be Bearer token")
)

func AuthMiddleware(tokenService *baseauth.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": ErrTokenMissing.Error(),
				"data":    nil,
			})
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": ErrTokenMalformed.Error(),
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		claims, err := tokenService.ParseAccessToken(strings.TrimSpace(parts[1]))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": fmt.Sprintf("invalid token: %s", err.Error()),
				"data":    nil,
			})
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
