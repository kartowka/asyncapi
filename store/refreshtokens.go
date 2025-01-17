package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/antfley/asyncapi/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenStore struct {
	repo *repository.Queries
}

func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		repo: repository.New(db),
	}
}
func (s *RefreshTokenStore) getBase64HashFromToken(token *jwt.Token) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(token.Raw))
	hashedToken := hash.Sum(nil)
	base64TokenHash := base64.StdEncoding.EncodeToString(hashedToken)
	return base64TokenHash, nil
}
func (s *RefreshTokenStore) Create(ctx context.Context, userID uint, token *jwt.Token) (*repository.RefreshToken, error) {
	base64TokenHash, err := s.getBase64HashFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get base64 hash from token: %w", err)
	}
	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration time: %w", err)
	}
	refreshToken := repository.CreateRefreshTokenParams{
		UserID:      userID,
		HashedToken: base64TokenHash,
		ExpiresAt:   expiresAt.Time,
	}
	result, err := s.repo.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	rt, err := s.repo.GetRefreshTokenByID(ctx, uint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token by id: %w", err)
	}
	return &rt, nil
}

func (s *RefreshTokenStore) DeleteUserTokens(ctx context.Context, id uint) (sql.Result, error) {
	result, err := s.repo.DeleteUserTokens(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return result, nil
}
