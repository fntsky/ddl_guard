package jwt

import (
	"errors"
	"fmt"
	"strings"

	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	gjwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = apperrors.ErrInvalidToken
	ErrTokenExpired     = apperrors.ErrTokenExpired
	ErrInvalidSignature = apperrors.New(401, "invalid signature")
	ErrInvalidClaims    = apperrors.New(401, "invalid claims")
)

type Claims struct {
	UserUUID string `json:"user_uuid"`
	TokenUse string `json:"token_use"`
	TokenID  string `json:"jti,omitempty"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

func GenerateToken(secret string, claims Claims) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", fmt.Errorf("secret is empty")
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
		return nil, fmt.Errorf("secret is empty")
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
