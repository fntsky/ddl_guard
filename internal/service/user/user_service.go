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
	ExistsByAuthIdentifier(ctx context.Context, authType string, authIdentifier string) (bool, error)
	CreateUserWithAuth(ctx context.Context, user *entity.User, auth *entity.UserAuth) error
	GetUserWithAuthByIdentifier(ctx context.Context, authType string, authIdentifier string) (*entity.User, *entity.UserAuth, error)
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

	exists, err := s.repo.ExistsByAuthIdentifier(ctx, entity.UserAuthTypeEmail, email)
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
		UUID:      uuid.GenerateUUID(),
		Username:  strings.TrimSpace(req.Username),
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	authInfo := &entity.UserAuth{
		AuthType:       entity.UserAuthTypeEmail,
		AuthIdentifier: email,
		CredentialHash: string(pwdHash),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err = s.repo.CreateUserWithAuth(ctx, user, authInfo); err != nil {
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

	user, auth, err := s.repo.GetUserWithAuthByIdentifier(ctx, entity.UserAuthTypeEmail, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(auth.CredentialHash), []byte(req.Password))
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
