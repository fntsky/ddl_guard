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
