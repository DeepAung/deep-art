package mytoken

import (
	"errors"
	"fmt"
	"time"

	"github.com/DeepAung/deep-art/pkg/config"
	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenType struct {
	Subject   string
	ExpiresAt func(cfg *config.JwtConfig) time.Time
}

var (
	Access = TokenType{
		Subject: "access-token",
		ExpiresAt: func(cfg *config.JwtConfig) time.Time {
			return time.Now().Add(cfg.AccessExpires)
		},
	}

	Refresh = TokenType{
		Subject: "refresh-token",
		ExpiresAt: func(cfg *config.JwtConfig) time.Time {
			return time.Now().Add(cfg.RefreshExpires)
		},
	}
)

type JwtClaims struct {
	Payload Payload
	jwt.RegisteredClaims
}

type Payload struct {
	UserId  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
}

func GenerateToken(
	cfg *config.JwtConfig,
	tokenType TokenType,
	payload Payload,
) (string, error) {

	claims := &JwtClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "deep-art-api",
			Subject:   tokenType.Subject,
			Audience:  []string{"users", "admin"},
			ExpiresAt: jwt.NewNumericDate(tokenType.ExpiresAt(cfg)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(cfg.SecretKey)
	if err != nil {
		return "", fmt.Errorf("generate token failed: %v", err)
	}

	return tokenString, nil
}

func ParseToken(
	cfg *config.JwtConfig,
	tokenType TokenType,
	tokenString string,
) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, keyFunc([]byte(cfg.SecretKey)))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); !ok {
		return nil, errors.New("invalid claims type")
	} else {
		return claims, nil
	}
}

func VerifyToken(cfg *config.JwtConfig, tokenType TokenType, tokenString string) error {
	_, err := ParseToken(cfg, tokenType, tokenString)
	return err
}

func keyFunc(key []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return key, nil
	}
}
