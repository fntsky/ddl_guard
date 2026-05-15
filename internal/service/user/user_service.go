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
	"github.com/fntsky/ddl_guard/internal/service/wechat"
	"github.com/fntsky/ddl_guard/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error
	GetUserByID(ctx context.Context, userID int64) (*entity.User, error)
	GetUserEmailsByIDs(ctx context.Context, userIDs []int64) (map[int64]string, error)
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

type UserAuthRepo interface {
	GetByTypeAndIdentifier(ctx context.Context, authType, identifier string) (*entity.UserAuth, error)
	Create(ctx context.Context, auth *entity.UserAuth) error
}

type UserService struct {
	repo         UserRepo
	userAuthRepo UserAuthRepo
	emailOTP     otp.OTP
	authService  *authsvc.AuthService
	wechatSvc    *wechat.WechatService
}

func NewUserService(repo UserRepo, userAuthRepo UserAuthRepo, emailOTP otp.OTP, authService *authsvc.AuthService, wechatSvc *wechat.WechatService) *UserService {
	return &UserService{
		repo:         repo,
		userAuthRepo: userAuthRepo,
		emailOTP:     emailOTP,
		authService:  authService,
		wechatSvc:    wechatSvc,
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
		Email:        entity.StrPtr(email),
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
	VerificationTypePhone = "phone"
)

// SendPasswordResetCode 发送密码重置验证码
func (s *UserService) SendPasswordResetCode(ctx context.Context, req *schema.SendPasswordResetCodeReq) error {
	target := strings.TrimSpace(req.Target)

	switch req.Type {
	case VerificationTypeEmail:
		target = normalizeEmail(target)
	case VerificationTypePhone:
		// TODO: 实现短信发送
		return apperrors.ErrSMSServiceDisabled
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
	case VerificationTypePhone:
		// TODO: 实现短信验证码验证
		return apperrors.ErrSMSServiceDisabled
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
	case VerificationTypePhone:
		user, err = s.repo.GetUserByPhone(ctx, target)
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

// ========== 手机号相关方法（空实现） ==========

// SendPhoneVerificationCode 发送手机注册验证码
func (s *UserService) SendPhoneVerificationCode(ctx context.Context, req *schema.SendPhoneVerificationCodeReq) error {
	// TODO: 实现短信发送
	return apperrors.ErrSMSServiceDisabled
}

// SendPhoneLoginCode 发送手机登录验证码
func (s *UserService) SendPhoneLoginCode(ctx context.Context, req *schema.SendPhoneLoginCodeReq) error {
	// TODO: 实现短信发送
	return apperrors.ErrSMSServiceDisabled
}

// RegisterByPhone 手机号注册
func (s *UserService) RegisterByPhone(ctx context.Context, req *schema.RegisterUserByPhoneReq) (*schema.RegisterUserByPhoneResp, error) {
	// TODO: 实现手机号注册
	return nil, apperrors.ErrSMSServiceDisabled
}

// LoginByPhone 手机号+密码登录
func (s *UserService) LoginByPhone(ctx context.Context, req *schema.LoginByPhoneReq) (*schema.LoginByPhoneResp, error) {
	phone := strings.TrimSpace(req.Phone)

	user, err := s.repo.GetUserByPhone(ctx, phone)
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

	return &schema.LoginByPhoneResp{
		UUID:         user.UUID,
		Username:     user.Username,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// LoginByPhoneCode 手机号+验证码登录
func (s *UserService) LoginByPhoneCode(ctx context.Context, req *schema.LoginByPhoneCodeReq) (*schema.LoginByPhoneCodeResp, error) {
	// TODO: 实现手机号验证码登录
	return nil, apperrors.ErrSMSServiceDisabled
}

// LoginByWechat 微信登录（首次自动注册）
func (s *UserService) LoginByWechat(ctx context.Context, req *schema.LoginByWechatReq) (*schema.LoginByWechatResp, error) {
	result, err := s.wechatSvc.Code2Session(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	userAuth, err := s.userAuthRepo.GetByTypeAndIdentifier(ctx, entity.UserAuthTypeWechat, result.OpenID)
	if err != nil {
		return nil, err
	}

	var user *entity.User

	if userAuth != nil {
		user, err = s.repo.GetUserByID(ctx, userAuth.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, apperrors.ErrUserNotFound
		}
	} else {
		now := time.Now()
		user = &entity.User{
			UUID:         uuid.GenerateUUID(),
			Username:     "wx_" + result.OpenID[:min(8, len(result.OpenID))],
			PasswordHash: string(mustGenerateRandomHash()),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err = s.repo.CreateUser(ctx, user); err != nil {
			return nil, err
		}

		authRecord := &entity.UserAuth{
			UserID:         user.ID,
			AuthType:       entity.UserAuthTypeWechat,
			AuthIdentifier: result.OpenID,
			AuthMeta:       result.SessionKey,
		}
		if err = s.userAuthRepo.Create(ctx, authRecord); err != nil {
			return nil, err
		}
	}

	tokenPair, err := s.authService.IssueTokensForUser(ctx, user.ID, user.UUID)
	if err != nil {
		return nil, err
	}

	return &schema.LoginByWechatResp{
		UUID:         user.UUID,
		Username:     user.Username,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func mustGenerateRandomHash() []byte {
	pwd, _ := bcrypt.GenerateFromPassword([]byte(uuid.GenerateUUID()), bcrypt.DefaultCost)
	return pwd
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
