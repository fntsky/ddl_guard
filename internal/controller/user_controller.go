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
// @Router /api/v1/users/email/verification-codes [post]
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
// @Router /api/v1/users/registrations/email [post]
func (uc *UserController) RegisterUserByEmail(ctx *gin.Context) {
	req := &schema.RegisterUserByEmailReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := uc.userService.RegisterByEmail(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// @Summary 邮箱登录
// @Description 使用邮箱和密码登录
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.LoginByEmailReq true "Login By Email Request"
// @success 200 {object} handler.resp{data=schema.LoginByEmailResp} "success"
// @Router /api/v1/users/sessions/email [post]
func (uc *UserController) LoginByEmail(ctx *gin.Context) {
	req := &schema.LoginByEmailReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := uc.userService.LoginByEmail(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// SendPasswordResetCode 发送密码重置验证码
// @Summary 发送密码重置验证码
// @Description 根据指定的验证类型发送验证码到目标地址。目前支持 email 类型，验证码有效期为 5 分钟。
// @Description type: 验证方式，目前支持 "email"
// @Description target: 验证目标，邮箱类型时为邮箱地址
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.SendPasswordResetCodeReq true "发送密码重置验证码请求"
// @success 200 {object} handler.resp "验证码发送成功"
// @failure 400 {object} handler.resp "请求参数错误"
// @failure 500 {object} handler.resp "服务器内部错误"
// @Router /api/v1/users/password/reset-codes [post]
func (uc *UserController) SendPasswordResetCode(ctx *gin.Context) {
	req := &schema.SendPasswordResetCodeReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := uc.userService.SendPasswordResetCode(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// ChangePassword 修改密码
// @Summary 通过验证码修改密码
// @Description 使用验证码验证身份后修改用户密码。验证码通过 SendPasswordResetCode 接口发送。
// @Description type: 验证方式，需与发送验证码时的类型一致，目前支持 "email"
// @Description target: 验证目标，需与发送验证码时的目标一致
// @Description code: 收到的验证码，6 位数字
// @Description new_password: 新密码，最少 6 个字符
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.ChangePasswordReq true "修改密码请求"
// @success 200 {object} handler.resp "密码修改成功"
// @failure 400 {object} handler.resp "请求参数错误或验证码无效"
// @failure 404 {object} handler.resp "用户不存在"
// @failure 500 {object} handler.resp "服务器内部错误"
// @Router /api/v1/users/password [put]
func (uc *UserController) ChangePassword(ctx *gin.Context) {
	req := &schema.ChangePasswordReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := uc.userService.ChangePassword(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
