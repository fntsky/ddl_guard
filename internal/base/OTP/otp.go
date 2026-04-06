package otp

import "context"

/*
OTP(One-Time Password)验证码
*/

// Purpose constants for different OTP use cases
const (
	PurposeRegister       = "register"        // 用户注册
	PurposeLogin          = "login"           // 用户登录
	PurposeResetPassword  = "reset_password"  // 重置密码
	PurposeBindEmail      = "bind_email"      // 绑定邮箱
	PurposeVerifyEmail    = "verify_email"    // 验证邮箱
	PurposeChangeEmail    = "change_email"    // 更换邮箱
)

type OTP interface {
	Send(ctx context.Context, purpose string, target string) error
	Verify(ctx context.Context, purpose string, target string, code string) (bool, error)
}
