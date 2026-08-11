package tokenprovider

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nemolucky/patisserie-go/internal/config"
	"github.com/nemolucky/patisserie-go/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider struct{
	AccessSecret string
	AccessExpire time.Duration
	RefreshSecret string
	RefreshExpire time.Duration
}

func NewTokenProvider(cfg config.JWTConfig) *TokenProvider{
	return &TokenProvider{
		AccessSecret: cfg.AccessSecret,
		AccessExpire: cfg.AccessExpiration,
		RefreshSecret: cfg.RefreshSecret,
		RefreshExpire: cfg.RefreshExpiration,
	}
}

func (t *TokenProvider) generateAccessToken(userID uuid.UUID) (string, error) {
	claims := domain.Claims{
		UserID: userID,
		TokenType: domain.AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.AccessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.AccessSecret))
}


func (t *TokenProvider) generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := domain.Claims{
		UserID: userID,
		TokenType: domain.RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.RefreshSecret))
}

func (t *TokenProvider) GeneratePairTokens(userID uuid.UUID) (*domain.TokenPair, error) {
	access, err := t.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refresh, err := t.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (t *TokenProvider) Parse(tokenString string, tokenType domain.TokenType) (*domain.Claims, error) {

	var tokenSecret string
	switch tokenType {
	case domain.AccessTokenType:
		tokenSecret = t.AccessSecret
	case domain.RefreshTokenType:
		tokenSecret = t.RefreshSecret
	default:
		return nil, errors.New("invalid token type")
	}

	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(tokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != tokenType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}