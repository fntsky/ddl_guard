package jwt

import (
	"errors"
	"net/http"
	"strings"

	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	gjwt "github.com/golang-jwt/jwt/v5"
)

var (
	// 直接使用 apperrors 中的错误定义
	ErrInvalidToken     = apperrors.ErrInvalidToken
	ErrTokenExpired     = apperrors.ErrTokenExpired
	ErrInvalidSignature = apperrors.New(http.StatusUnauthorized, apperrors.CodeInvalidToken, "invalid signature")
	ErrInvalidClaims    = apperrors.New(http.StatusUnauthorized, apperrors.CodeInvalidToken, "invalid claims")
)

type Claims struct {
	UserUUID string `json:"user_uuid"`
	TokenUse string `json:"token_use"`
	TokenID  string `json:"jti,omitempty"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

func (c *Claims) GetUserUUID() string {
	return c.UserUUID
}

func GenerateToken(secret string, claims Claims) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", apperrors.ErrTokenConfigInvalid
	}
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, gjwt.MapClaims{
		"user_uuid": claims.UserUUID,
		"token_use": claims.TokenUse,
		"jti":       claims.TokenID,
		"exp":       claims.Exp,
		"iat":       claims.Iat,
	})
	return token.SignedString([]byte(secret))
}

func ParseToken(secret string, token string) (*Claims, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, apperrors.ErrTokenConfigInvalid
	}

	parsedToken, err := gjwt.Parse(token, func(t *gjwt.Token) (any, error) {
		if t.Method == nil || t.Method.Alg() != gjwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}, gjwt.WithValidMethods([]string{gjwt.SigningMethodHS256.Alg()}))
	if err != nil {
		if errors.Is(err, gjwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, gjwt.ErrTokenSignatureInvalid) {
			return nil, ErrInvalidSignature
		}
		return nil, ErrInvalidToken
	}
	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	claimsMap, ok := parsedToken.Claims.(gjwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	exp, err := claimsMap.GetExpirationTime()
	if err != nil || exp == nil {
		return nil, ErrInvalidClaims
	}
	iat, err := claimsMap.GetIssuedAt()
	if err != nil || iat == nil {
		return nil, ErrInvalidClaims
	}

	userUUID, ok := claimsMap["user_uuid"].(string)
	if !ok {
		return nil, ErrInvalidClaims
	}
	tokenUse, ok := claimsMap["token_use"].(string)
	if !ok {
		return nil, ErrInvalidClaims
	}
	tokenID, _ := claimsMap["jti"].(string)

	return &Claims{
		UserUUID: userUUID,
		TokenUse: tokenUse,
		TokenID:  tokenID,
		Exp:      exp.Unix(),
		Iat:      iat.Unix(),
	}, nil
}