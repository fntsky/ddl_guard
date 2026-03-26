package controller

import (
	"github.com/fntsky/ddl_guard/internal/base/handler"
	"github.com/fntsky/ddl_guard/internal/schema"
	"github.com/fntsky/ddl_guard/internal/service/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *user.UserService
}

func NewUserController(userService *user.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码用于注册
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.SendEmailVerificationCodeReq true "Send Email Verification Code Request"
// @success 200 {object} handler.resp "success"
// @Router /users/email/verification-codes [post]
func (uc *UserController) SendEmailVerificationCode(ctx *gin.Context) {
	req := &schema.SendEmailVerificationCodeReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := uc.userService.SendEmailVerificationCode(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// @Summary 邮箱注册
// @Description 使用邮箱验证码注册用户
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.RegisterUserByEmailReq true "Register User By Email Request"
// @success 200 {object} handler.resp{data=schema.RegisterUserByEmailResp} "success"
// @Router /users/registrations/email [post]
func (uc *UserController) RegisterUserByEmail(ctx *gin.Context) {
	req := &schema.RegisterUserByEmailReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := uc.userService.RegisterByEmail(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}
