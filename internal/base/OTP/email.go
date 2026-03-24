package otp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/conf"
)

const (
	defaultCodeLength = 6
	defaultCodeTTL    = 5 * time.Minute
	defaultSendTTL    = 1 * time.Minute
	defaultFromName   = "ddl_guard"
)

type otpRepo interface {
	storeCode(target string, code string, expiresAt time.Time) error
	getRateLimit(target string) (canSend bool, retryAfter time.Duration, err error)
	getCode(target string) (code string, expiresAt time.Time, found bool, err error)
}

type SMTPEmailOTP struct {
	host     string
	port     int
	username string
	password string
	codeTTL  time.Duration
	sendTTL  time.Duration
	codes    otpRepo
}

func (s *SMTPEmailOTP) Send(ctx context.Context, target string) error {
	if err := validateEmail(target); err != nil {
		return err
	}
	//检查发送频率
	canSend, retryAfter, err := s.codes.getRateLimit(target)
	if err != nil {
		return fmt.Errorf("check rate limit: %w", err)
	}
	if !canSend {
		return fmt.Errorf("rate limit exceeded, please try again after %v", retryAfter)
	}
	code, err := GenerateNumericCode(defaultCodeLength)
	if err != nil {
		return fmt.Errorf("generate otp code: %w", err)
	}
	subject := "您的验证码"
	body := fmt.Sprintf("您的验证码是: %s\n此验证码将在 %d 分钟后过期。", code, int(s.codeTTL.Minutes()))

	if err = s.sendMail(ctx, target, subject, body); err != nil {
		return err
	}

	s.codes.storeCode(target, code, time.Now().Add(s.codeTTL))
	return nil
}

func (s *SMTPEmailOTP) Verify(ctx context.Context, target string, code string) (bool, error) {
	if err := validateEmail(target); err != nil {
		return false, err
	}
	storedCode, expiresAt, found, err := s.codes.getCode(target)
	if err != nil {
		return false, fmt.Errorf("get stored code: %w", err)
	}
	if !found || time.Now().After(expiresAt) {
		return false, nil // code not found or expired
	}
	return storedCode == code, nil
}

func NewSMTPEmailOTP() *SMTPEmailOTP {
	cfg := conf.Global()
	if cfg == nil || cfg.EMAIL_OTP.SMTP.Host == "" {
		panic("EMAIL_OTP SMTP config is not loaded or empty")
	}
	smtpCfg := cfg.EMAIL_OTP.SMTP
	return &SMTPEmailOTP{
		host:     smtpCfg.Host,
		port:     smtpCfg.Port,
		username: smtpCfg.Username,
		password: smtpCfg.Password,
		codeTTL:  defaultCodeTTL,
		sendTTL:  defaultSendTTL,
		codes:    nil, //TODO: 实现一个redis的otpRepo
	}
}

func (s *SMTPEmailOTP) sendMail(ctx context.Context, target string, subject string, body string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	msg := buildMailMessage(s.username, target, subject, body)
	client, err := s.newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("AUTH"); ok {
		auth := smtp.PlainAuth("", s.username, s.password, s.host)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
	}
	if err = client.Mail(s.username); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}
	if err = client.Rcpt(target); err != nil {
		return fmt.Errorf("smtp RCPT TO failed: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}
	if _, err = w.Write([]byte(msg)); err != nil {
		_ = w.Close()
		return fmt.Errorf("write mail body failed: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("close mail body failed: %w", err)
	}
	if err = client.Quit(); err != nil {
		return fmt.Errorf("smtp quit failed: %w", err)
	}
	return nil
}

func (s *SMTPEmailOTP) newClient() (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	if s.port == 465 {
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			ServerName: s.host,
			MinVersion: tls.VersionTLS12,
		})
		if err != nil {
			return nil, fmt.Errorf("tls dial smtp server failed: %w", err)
		}
		client, err := smtp.NewClient(conn, s.host)
		if err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("create smtp client failed: %w", err)
		}
		return client, nil
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("dial smtp server failed: %w", err)
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(&tls.Config{
			ServerName: s.host,
			MinVersion: tls.VersionTLS12,
		}); err != nil {
			_ = client.Close()
			return nil, fmt.Errorf("starttls failed: %w", err)
		}
	}
	return client, nil
}

func buildMailMessage(from string, to string, subject string, body string) string {
	return fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n",
		defaultFromName,
		from,
		to,
		subject,
		body,
	)
}

func validateEmail(target string) error {
	if _, err := mail.ParseAddress(target); err != nil {
		return fmt.Errorf("invalid target email: %w", err)
	}
	return nil
}
