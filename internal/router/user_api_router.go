package router

import (
	"github.com/fntsky/ddl_guard/internal/controller"
	"github.com/gin-gonic/gin"
)

type UserApiRouter struct {
	userController *controller.UserController
}

func NewUserApiRouter(userController *controller.UserController) *UserApiRouter {
	return &UserApiRouter{userController: userController}
}

func (a *UserApiRouter) Register(r *gin.RouterGroup) {
	userGroup := r.Group("/users")
	userGroup.POST("/email/verification-codes", a.userController.SendEmailVerificationCode)
	userGroup.POST("/registrations/email", a.userController.RegisterUserByEmail)
	userGroup.POST("/sessions/email", a.userController.LoginByEmail)
}
