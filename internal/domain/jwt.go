package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

var(
	AccessTokenType TokenType = "access" 
	RefreshTokenType TokenType = "refresh" 
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	TokenType TokenType 	`json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}