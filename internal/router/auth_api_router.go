package router

import (
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/gin-gonic/gin"
)

type AuthApiRouter struct {
	authController *controller.AuthController
}

func NewAuthApiRouter(authController *controller.AuthController) *AuthApiRouter {
	return &AuthApiRouter{authController: authController}
}

func (a *AuthApiRouter) Register(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	authGroup.POST("/refresh-tokens", a.authController.RefreshToken)
}
