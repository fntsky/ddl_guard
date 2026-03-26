package otp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/fntsky/ddl_guard/internal/base/email"
)

const (
	defaultCodeLength = 6
	defaultCodeTTL    = 5 * time.Minute
)

type otpRepo interface {
	StoreCode(target string, code string, expiresAt time.Time) error
	GetCode(target string) (code string, expiresAt time.Time, found bool, err error)
	DeleteCode(target string) error
}

var (
	ErrCodeStoreNotConfigured = errors.New("otp code store is not configured")
	ErrEmailOTPDisabled       = errors.New("email otp is disabled")
)

type EmailOTP struct {
	sender   email.Sender
	codeTTL  time.Duration
	codeRepo otpRepo
}

type disabledOTP struct{}

func (d *disabledOTP) Send(ctx context.Context, target string) error {
	return ErrEmailOTPDisabled
}

func (d *disabledOTP) Verify(ctx context.Context, target string, code string) (bool, error) {
	return false, ErrEmailOTPDisabled
}

func (s *EmailOTP) Send(ctx context.Context, target string) error {
	target = strings.TrimSpace(target)
	code, err := GenerateNumericCode(defaultCodeLength)
	if err != nil {
		return fmt.Errorf("generate otp code: %w", err)
	}
	if err = s.sender.Send(ctx, target, "您的验证码", fmt.Sprintf("您的验证码是: %s\n验证码将在 %d 分钟后过期。", code, int(s.codeTTL.Minutes()))); err != nil {
		return err
	}
	if s.codeRepo == nil {
		return nil
	}
	if err = s.codeRepo.StoreCode(target, code, time.Now().Add(s.codeTTL)); err != nil {
		return fmt.Errorf("store otp code failed: %w", err)
	}
	return nil
}

func (s *EmailOTP) Verify(ctx context.Context, target string, code string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	if s.codeRepo == nil {
		return false, ErrCodeStoreNotConfigured
	}
	target = strings.TrimSpace(target)
	code = strings.TrimSpace(code)
	storedCode, expiresAt, found, err := s.codeRepo.GetCode(target)
	if err != nil {
		return false, fmt.Errorf("get stored code: %w", err)
	}
	if !found || time.Now().After(expiresAt) || code == "" || storedCode != code {
		return false, nil
	}
	if err = s.codeRepo.DeleteCode(target); err != nil {
		return false, fmt.Errorf("delete used code failed: %w", err)
	}
	return true, nil
}

func NewSMTPEmailOTP() OTP {
	cfg := conf.Global()
	if cfg == nil {
		return &disabledOTP{}
	}
	smtpCfg := cfg.EMAIL_OTP.SMTP
	if strings.TrimSpace(smtpCfg.Host) == "" || smtpCfg.Port == 0 ||
		strings.TrimSpace(smtpCfg.Username) == "" || strings.TrimSpace(smtpCfg.Password) == "" {
		return &disabledOTP{}
	}
	sender := email.NewSMTPSender(smtpCfg.Host, smtpCfg.Port, smtpCfg.Username, smtpCfg.Password)
	return &EmailOTP{
		sender:   sender,
		codeTTL:  defaultCodeTTL,
		codeRepo: nil,
	}
}
