package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

const defaultFromName = "ddl_guard"

type Sender interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
}

func NewSMTPSender(host string, port int, username string, password string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *SMTPSender) Send(ctx context.Context, to string, subject string, body string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	msg := buildMailMessage(s.username, to, subject, body)
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
	if err = client.Rcpt(to); err != nil {
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

func (s *SMTPSender) newClient() (*smtp.Client, error) {
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
