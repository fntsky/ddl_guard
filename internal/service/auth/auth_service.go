package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	baseauth "github.com/fntsky/ddl_guard/internal/base/auth"
	"github.com/fntsky/ddl_guard/internal/entity"
	"github.com/fntsky/ddl_guard/internal/schema"
	pkgjwt "github.com/fntsky/ddl_guard/pkg/jwt"
	"github.com/fntsky/ddl_guard/pkg/uuid"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrRefreshTokenRevoked = errors.New("refresh token revoked")
	ErrSessionNotFound     = errors.New("session not found")
)

type SessionRepo interface {
	CreateSession(ctx context.Context, session *entity.UserSession) error
	GetByTokenID(ctx context.Context, tokenID string) (*entity.UserSession, bool, error)
	RotateSession(ctx context.Context, currentTokenID string, newSession *entity.UserSession) error
	RevokeByTokenID(ctx context.Context, tokenID string) error
}

type AuthService struct {
	tokenService *baseauth.TokenService
	sessionRepo  SessionRepo
}

func NewAuthService(tokenService *baseauth.TokenService, sessionRepo SessionRepo) *AuthService {
	return &AuthService{
		tokenService: tokenService,
		sessionRepo:  sessionRepo,
	}
}

func (s *AuthService) IssueTokensForUser(ctx context.Context, userID int64, userUUID string) (*schema.TokenPairResp, error) {
	tokenID := uuid.GenerateUUID()
	refreshToken, err := s.tokenService.GenerateRefreshToken(userUUID, tokenID)
	if err != nil {
		return nil, err
	}
	refreshClaims, err := s.tokenService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.tokenService.GenerateAccessToken(userUUID)
	if err != nil {
		return nil, err
	}

	session := &entity.UserSession{
		UserID:            userID,
		TokenID:           tokenID,
		RefreshTokenHash:  hashRefreshToken(refreshToken),
		ExpiresAt:         time.Unix(refreshClaims.Exp, 0),
		ReplacedByTokenID: "",
	}
	if err = s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &schema.TokenPairResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *schema.RefreshTokenReq) (*schema.TokenPairResp, error) {
	rawToken := strings.TrimSpace(req.RefreshToken)
	if rawToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	claims, err := s.tokenService.ParseRefreshToken(rawToken)
	if err != nil {
		if errors.Is(err, pkgjwt.ErrTokenExpired) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, ErrInvalidRefreshToken
	}
	if strings.TrimSpace(claims.TokenID) == "" {
		return nil, ErrInvalidRefreshToken
	}

	session, has, err := s.sessionRepo.GetByTokenID(ctx, claims.TokenID)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrSessionNotFound
	}
	if session.RevokedAt != nil {
		return nil, ErrRefreshTokenRevoked
	}
	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.RevokeByTokenID(ctx, claims.TokenID)
		return nil, ErrRefreshTokenExpired
	}
	if session.RefreshTokenHash != hashRefreshToken(rawToken) {
		return nil, ErrInvalidRefreshToken
	}

	newTokenID := uuid.GenerateUUID()
	newRefreshToken, err := s.tokenService.GenerateRefreshToken(claims.UserUUID, newTokenID)
	if err != nil {
		return nil, err
	}
	newRefreshClaims, err := s.tokenService.ParseRefreshToken(newRefreshToken)
	if err != nil {
		return nil, err
	}
	newAccessToken, err := s.tokenService.GenerateAccessToken(claims.UserUUID)
	if err != nil {
		return nil, err
	}

	newSession := &entity.UserSession{
		UserID:            session.UserID,
		TokenID:           newTokenID,
		RefreshTokenHash:  hashRefreshToken(newRefreshToken),
		ExpiresAt:         time.Unix(newRefreshClaims.Exp, 0),
		ReplacedByTokenID: "",
	}
	if err = s.sessionRepo.RotateSession(ctx, claims.TokenID, newSession); err != nil {
		return nil, fmt.Errorf("rotate refresh token failed: %w", err)
	}

	return &schema.TokenPairResp{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
