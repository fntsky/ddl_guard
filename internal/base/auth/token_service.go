package auth

import (
	"strings"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	pkgjwt "github.com/fntsky/ddl_guard/pkg/jwt"
)

const (
	TokenUseAccess  = "access"
	TokenUseRefresh = "refresh"
)

type TokenService struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService() (*TokenService, error) {
	cfg := conf.Global()
	if cfg == nil {
		return nil, apperrors.ErrTokenConfigInvalid
	}
	j := cfg.JWT
	if strings.TrimSpace(j.Secret) == "" {
		return nil, apperrors.ErrTokenConfigInvalid
	}

	accessTTL := time.Duration(j.AccessTTLMinutes) * time.Minute
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	refreshTTL := time.Duration(j.RefreshTTLHours) * time.Hour
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}

	return &TokenService{
		secret:     j.Secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

func (s *TokenService) GenerateTokenPair(userUUID string) (accessToken string, refreshToken string, err error) {
	refreshTokenID := ""
	return s.GenerateTokenPairWithRefreshID(userUUID, refreshTokenID)
}

func (s *TokenService) GenerateTokenPairWithRefreshID(userUUID string, refreshTokenID string) (accessToken string, refreshToken string, err error) {
	accessToken, err = s.GenerateAccessToken(userUUID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = s.GenerateRefreshToken(userUUID, refreshTokenID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *TokenService) GenerateAccessToken(userUUID string) (string, error) {
	now := time.Now().Unix()
	accessClaims := pkgjwt.Claims{
		UserUUID: userUUID,
		TokenUse: TokenUseAccess,
		Iat:      now,
		Exp:      now + int64(s.accessTTL.Seconds()),
	}
	return pkgjwt.GenerateToken(s.secret, accessClaims)
}

func (s *TokenService) GenerateRefreshToken(userUUID string, tokenID string) (string, error) {
	now := time.Now().Unix()
	refreshClaims := pkgjwt.Claims{
		UserUUID: userUUID,
		TokenUse: TokenUseRefresh,
		TokenID:  tokenID,
		Iat:      now,
		Exp:      now + int64(s.refreshTTL.Seconds()),
	}
	return pkgjwt.GenerateToken(s.secret, refreshClaims)
}

func (s *TokenService) ParseAccessToken(token string) (*pkgjwt.Claims, error) {
	claims, err := pkgjwt.ParseToken(s.secret, token)
	if err != nil {
		return nil, err
	}
	if claims.TokenUse != TokenUseAccess {
		return nil, pkgjwt.ErrInvalidClaims
	}
	return claims, nil
}

func (s *TokenService) ParseRefreshToken(token string) (*pkgjwt.Claims, error) {
	claims, err := pkgjwt.ParseToken(s.secret, token)
	if err != nil {
		return nil, err
	}
	if claims.TokenUse != TokenUseRefresh {
		return nil, pkgjwt.ErrInvalidClaims
	}
	return claims, nil
}
