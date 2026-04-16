package controller

import (
	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/schema"
	authsvc "github.com/fntsky/ddl_guard/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *authsvc.AuthService
}

func NewAuthController(authService *authsvc.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// @Summary 刷新Token
// @Description 使用refresh token换取新的access token和refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param req body schema.RefreshTokenReq true "Refresh Token Request"
// @success 200 {object} handler.resp{data=schema.TokenPairResp} "success"
// @Router /api/v1/auth/refresh-tokens [post]
func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	req := &schema.RefreshTokenReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := ac.authService.RefreshToken(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}
