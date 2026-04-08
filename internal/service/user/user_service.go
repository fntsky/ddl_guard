package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/OTP"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/schema"
	authsvc "github.com/fntsky/ddl_guard/internal/service/auth"
	"github.com/fntsky/ddl_guard/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error
	GetUserByID(ctx context.Context, userID int64) (*entity.User, error)
}

type UserService struct {
	repo        UserRepo
	emailOTP    otp.OTP
	authService *authsvc.AuthService
}

func NewUserService(repo UserRepo, emailOTP otp.OTP, authService *authsvc.AuthService) *UserService {
	return &UserService{
		repo:        repo,
		emailOTP:    emailOTP,
		authService: authService,
	}
}

func (s *UserService) SendEmailVerificationCode(ctx context.Context, req *schema.SendEmailVerificationCodeReq) error {
	err := s.emailOTP.Send(ctx, otp.PurposeRegister, normalizeEmail(req.Email))
	if errors.Is(err, apperrors.ErrEmailOTPDisabled) {
		return apperrors.ErrEmailOTPDisabled
	}
	return err
}

func (s *UserService) RegisterByEmail(ctx context.Context, req *schema.RegisterUserByEmailReq) (*schema.RegisterUserByEmailResp, error) {
	email := normalizeEmail(req.Email)
	ok, err := s.emailOTP.Verify(ctx, otp.PurposeRegister, email, strings.TrimSpace(req.Code))
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailOTPDisabled) {
			return nil, apperrors.ErrEmailOTPDisabled
		}
		if errors.Is(err, apperrors.ErrCodeStoreNotConfigured) {
			return nil, apperrors.ErrVerificationUnavailable
		}
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrInvalidVerificationCode
	}

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.ErrEmailAlreadyExists
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &entity.User{
		UUID:         uuid.GenerateUUID(),
		Username:     strings.TrimSpace(req.Username),
		Email:        email,
		PasswordHash: string(pwdHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err = s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	tokenPair, err := s.authService.IssueTokensForUser(ctx, user.ID, user.UUID)
	if err != nil {
		return nil, err
	}

	return &schema.RegisterUserByEmailResp{
		UUID:         user.UUID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *UserService) LoginByEmail(ctx context.Context, req *schema.LoginByEmailReq) (*schema.LoginByEmailResp, error) {
	email := normalizeEmail(req.Email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	tokenPair, err := s.authService.IssueTokensForUser(ctx, user.ID, user.UUID)
	if err != nil {
		return nil, err
	}

	return &schema.LoginByEmailResp{
		UUID:         user.UUID,
		Username:     user.Username,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

const (
	VerificationTypeEmail = "email"
)

// SendPasswordResetCode 发送密码重置验证码
func (s *UserService) SendPasswordResetCode(ctx context.Context, req *schema.SendPasswordResetCodeReq) error {
	target := strings.TrimSpace(req.Target)

	switch req.Type {
	case VerificationTypeEmail:
		target = normalizeEmail(target)
	default:
		return apperrors.ErrUnsupportedVerificationType
	}

	err := s.emailOTP.Send(ctx, otp.PurposeResetPassword, target)
	if errors.Is(err, apperrors.ErrEmailOTPDisabled) {
		return apperrors.ErrEmailOTPDisabled
	}
	return err
}

// ChangePassword 通过验证码修改密码
func (s *UserService) ChangePassword(ctx context.Context, req *schema.ChangePasswordReq) error {
	target := strings.TrimSpace(req.Target)

	switch req.Type {
	case VerificationTypeEmail:
		target = normalizeEmail(target)
	default:
		return apperrors.ErrUnsupportedVerificationType
	}

	// 验证验证码
	ok, err := s.emailOTP.Verify(ctx, otp.PurposeResetPassword, target, strings.TrimSpace(req.Code))
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailOTPDisabled) {
			return apperrors.ErrEmailOTPDisabled
		}
		if errors.Is(err, apperrors.ErrCodeStoreNotConfigured) {
			return apperrors.ErrVerificationUnavailable
		}
		return err
	}
	if !ok {
		return apperrors.ErrInvalidVerificationCode
	}

	// 获取用户
	var user *entity.User
	switch req.Type {
	case VerificationTypeEmail:
		user, err = s.repo.GetUserByEmail(ctx, target)
	default:
		return apperrors.ErrUnsupportedVerificationType
	}
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrUserNotFound
	}

	// 更新密码
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, user.ID, string(pwdHash))
}
