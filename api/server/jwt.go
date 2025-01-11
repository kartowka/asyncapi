package server

import (
	"fmt"
	"time"

	"github.com/antfley/asyncapi/config"
	"github.com/golang-jwt/jwt/v5"
)

var signingMethod = jwt.SigningMethodHS256

type JWTManager struct {
	config *config.Config
}

func NewJWTManager(config *config.Config) *JWTManager {
	return &JWTManager{config: config}
}

type TokenPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}
type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func (j *JWTManager) Parse(token string) (*jwt.Token, error) {
	parser := jwt.NewParser()
	jwtToken, err := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.config.JWT_SECRET), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return jwtToken, nil
}
func (j *JWTManager) IsAccessToken(token *jwt.Token) bool {
	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	if tokenType, ok := jwtClaims["token_type"]; ok {
		return tokenType == "access"
	}
	return false
}
func (j *JWTManager) GenerateTokenPair(userID uint) (*TokenPair, error) {
	now := time.Now()
	issuer := "http://localhost" + ":" + j.config.PORT
	jwtAccessToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	key := []byte(j.config.JWT_SECRET)
	var err error
	signedAccessToken, err := jwtAccessToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	jwtRefreshToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	accessToken, err := j.Parse(signedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}
	signedRefreshToken, err := jwtRefreshToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	refreshToken, err := j.Parse(signedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
