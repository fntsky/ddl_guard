package otp

import "context"

/*
OTP(One-Time Password)验证码
*/

type OTP interface {
	Send(ctx context.Context, target string) error
	Verify(ctx context.Context, target string, code string) (bool, error)
}
