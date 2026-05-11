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
	// 邮箱相关
	userGroup.POST("/email/verification-codes", a.userController.SendEmailVerificationCode)
	userGroup.POST("/registrations/email", a.userController.RegisterUserByEmail)
	userGroup.POST("/sessions/email", a.userController.LoginByEmail)
	// 手机号相关
	userGroup.POST("/phone/verification-codes", a.userController.SendPhoneVerificationCode)
	userGroup.POST("/phone/login-verification-codes", a.userController.SendPhoneLoginCode)
	userGroup.POST("/registrations/phone", a.userController.RegisterUserByPhone)
	userGroup.POST("/sessions/phone", a.userController.LoginByPhone)
	userGroup.POST("/sessions/phone/code", a.userController.LoginByPhoneCode)
	// 密码相关
	userGroup.POST("/password/reset-codes", a.userController.SendPasswordResetCode)
	userGroup.PUT("/password", a.userController.ChangePassword)
}
