package otp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/fntsky/ddl_guard/internal/base/data"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/base/email"
)

const (
	defaultCodeLength = 6
	defaultCodeTTL    = 5 * time.Minute
)

type otpRepo interface {
	StoreCode(purpose string, target string, code string, expiresAt time.Time) error
	GetCode(purpose string, target string) (code string, expiresAt time.Time, found bool, err error)
	DeleteCode(purpose string, target string) error
}

type EmailOTP struct {
	sender   email.Sender
	codeTTL  time.Duration
	codeRepo otpRepo
}

type disabledOTP struct{}

func (d *disabledOTP) Send(ctx context.Context, purpose string, target string) error {
	return apperrors.ErrEmailOTPDisabled
}

func (d *disabledOTP) Verify(ctx context.Context, purpose string, target string, code string) (bool, error) {
	return false, apperrors.ErrEmailOTPDisabled
}

func (s *EmailOTP) Send(ctx context.Context, purpose string, target string) error {
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
	if err = s.codeRepo.StoreCode(purpose, target, code, time.Now().Add(s.codeTTL)); err != nil {
		return fmt.Errorf("store otp code failed: %w", err)
	}
	return nil
}

func (s *EmailOTP) Verify(ctx context.Context, purpose string, target string, code string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	if s.codeRepo == nil {
		return false, apperrors.ErrCodeStoreNotConfigured
	}
	target = strings.TrimSpace(target)
	code = strings.TrimSpace(code)
	storedCode, expiresAt, found, err := s.codeRepo.GetCode(purpose, target)
	if err != nil {
		return false, fmt.Errorf("get stored code: %w", err)
	}
	if !found || time.Now().After(expiresAt) || code == "" || storedCode != code {
		return false, nil
	}
	if err = s.codeRepo.DeleteCode(purpose, target); err != nil {
		return false, fmt.Errorf("delete used code failed: %w", err)
	}
	return true, nil
}

func NewSMTPEmailOTP(redis *data.RedisClient) OTP {
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

	var codeRepo otpRepo
	if redis != nil {
		codeRepo = NewRedisOTPRepo(redis)
	}

	return &EmailOTP{
		sender:   sender,
		codeTTL:  defaultCodeTTL,
		codeRepo: codeRepo,
	}
}
