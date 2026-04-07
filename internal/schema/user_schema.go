package schema

/*
用户邮箱注册请求参数
用户名
邮箱
密码
验证码
*/
type SendEmailVerificationCodeReq struct {
	Email string `json:"email" binding:"required,email" example:"testuser@example.com"`
}

type RegisterUserByEmailReq struct {
	Username string `json:"username" binding:"required" example:"testuser"`
	Email    string `json:"email" binding:"required,email" example:"testuser@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Code     string `json:"code" binding:"required" example:"123456"`
}

/*
用户邮箱注册响应参数
UUID
refresh_token
access_token
*/
type RegisterUserByEmailResp struct {
	UUID         string `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

/*
邮箱登录请求参数
邮箱
密码
*/
type LoginByEmailReq struct {
	Email    string `json:"email" binding:"required,email" example:"testuser@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

/*
邮箱登录响应参数
UUID
用户名
refresh_token
access_token
*/
type LoginByEmailResp struct {
	UUID         string `json:"uuid" example:"7a178766-4b8e-4e99-ab4c-843f7dbd95fd"`
	Username     string `json:"username" example:"testuser"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

/*
发送密码重置验证码请求
type: 验证方式，目前支持 "email"
target: 验证目标，邮箱类型时为邮箱地址
*/
type SendPasswordResetCodeReq struct {
	Type   string `json:"type" binding:"required,oneof=email" example:"email"`
	Target string `json:"target" binding:"required" example:"user@example.com"`
}

/*
修改密码请求
type: 验证方式，需与发送验证码时的类型一致
target: 验证目标，需与发送验证码时的目标一致
code: 收到的验证码，6 位数字
new_password: 新密码，最少 6 个字符
*/
type ChangePasswordReq struct {
	Type        string `json:"type" binding:"required,oneof=email" example:"email"`
	Target      string `json:"target" binding:"required" example:"user@example.com"`
	Code        string `json:"code" binding:"required" example:"123456"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpassword123"`
}
