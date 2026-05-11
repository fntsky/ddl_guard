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
// @success 200 {object} handler.Response "success"
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
// @success 200 {object} handler.Response{data=schema.RegisterUserByEmailResp} "success"
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
// @success 200 {object} handler.Response{data=schema.LoginByEmailResp} "success"
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
// @Description 根据指定的验证类型发送验证码到目标地址。目前支持 email 和 phone 类型，验证码有效期为 5 分钟。
// @Description type: 验证方式，目前支持 "email", "phone"
// @Description target: 验证目标，邮箱类型时为邮箱地址，手机类型时为手机号
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.SendPasswordResetCodeReq true "发送密码重置验证码请求"
// @success 200 {object} handler.Response "验证码发送成功"
// @failure 400 {object} handler.Response "请求参数错误"
// @failure 500 {object} handler.Response "服务器内部错误"
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
// @Description type: 验证方式，需与发送验证码时的类型一致，目前支持 "email", "phone"
// @Description target: 验证目标，需与发送验证码时的目标一致
// @Description code: 收到的验证码，6 位数字
// @Description new_password: 新密码，最少 6 个字符
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.ChangePasswordReq true "修改密码请求"
// @success 200 {object} handler.Response "密码修改成功"
// @failure 400 {object} handler.Response "请求参数错误或验证码无效"
// @failure 404 {object} handler.Response "用户不存在"
// @failure 500 {object} handler.Response "服务器内部错误"
// @Router /api/v1/users/password [put]
func (uc *UserController) ChangePassword(ctx *gin.Context) {
	req := &schema.ChangePasswordReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := uc.userService.ChangePassword(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// ========== 手机号相关接口 ==========

// SendPhoneVerificationCode 发送手机注册验证码
// @Summary 发送手机注册验证码
// @Description 发送手机验证码用于注册（暂未实现）
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.SendPhoneVerificationCodeReq true "发送手机验证码请求"
// @success 200 {object} handler.Response "验证码发送成功"
// @failure 400 {object} handler.Response "请求参数错误"
// @failure 503 {object} handler.Response "短信服务未启用"
// @Router /api/v1/users/phone/verification-codes [post]
func (uc *UserController) SendPhoneVerificationCode(ctx *gin.Context) {
	req := &schema.SendPhoneVerificationCodeReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := uc.userService.SendPhoneVerificationCode(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// SendPhoneLoginCode 发送手机登录验证码
// @Summary 发送手机登录验证码
// @Description 发送手机验证码用于快捷登录（暂未实现）
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.SendPhoneLoginCodeReq true "发送手机登录验证码请求"
// @success 200 {object} handler.Response "验证码发送成功"
// @failure 400 {object} handler.Response "请求参数错误"
// @failure 503 {object} handler.Response "短信服务未启用"
// @Router /api/v1/users/phone/login-verification-codes [post]
func (uc *UserController) SendPhoneLoginCode(ctx *gin.Context) {
	req := &schema.SendPhoneLoginCodeReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := uc.userService.SendPhoneLoginCode(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// RegisterUserByPhone 手机号注册
// @Summary 手机号注册
// @Description 使用手机号验证码注册用户（暂未实现）
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.RegisterUserByPhoneReq true "手机号注册请求"
// @success 200 {object} handler.Response{data=schema.RegisterUserByPhoneResp} "注册成功"
// @failure 400 {object} handler.Response "请求参数错误或验证码无效"
// @failure 503 {object} handler.Response "短信服务未启用"
// @Router /api/v1/users/registrations/phone [post]
func (uc *UserController) RegisterUserByPhone(ctx *gin.Context) {
	req := &schema.RegisterUserByPhoneReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := uc.userService.RegisterByPhone(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// LoginByPhone 手机号+密码登录
// @Summary 手机号+密码登录
// @Description 使用手机号和密码登录
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.LoginByPhoneReq true "手机号登录请求"
// @success 200 {object} handler.Response{data=schema.LoginByPhoneResp} "登录成功"
// @failure 400 {object} handler.Response "请求参数错误"
// @failure 401 {object} handler.Response "手机号或密码错误"
// @Router /api/v1/users/sessions/phone [post]
func (uc *UserController) LoginByPhone(ctx *gin.Context) {
	req := &schema.LoginByPhoneReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := uc.userService.LoginByPhone(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// LoginByPhoneCode 手机号+验证码登录
// @Summary 手机号+验证码登录
// @Description 使用手机号和验证码快捷登录（暂未实现）
// @Tags User
// @Accept json
// @Produce json
// @Param req body schema.LoginByPhoneCodeReq true "手机号验证码登录请求"
// @success 200 {object} handler.Response{data=schema.LoginByPhoneCodeResp} "登录成功"
// @failure 400 {object} handler.Response "请求参数错误或验证码无效"
// @failure 503 {object} handler.Response "短信服务未启用"
// @Router /api/v1/users/sessions/phone/code [post]
func (uc *UserController) LoginByPhoneCode(ctx *gin.Context) {
	req := &schema.LoginByPhoneCodeReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	resp, err := uc.userService.LoginByPhoneCode(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}
