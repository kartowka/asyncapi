package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/antfley/asyncapi/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	repo *repository.Queries
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		repo: repository.New(db),
	}
}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*repository.User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)
	uuid := uuid.New()
	user := repository.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPasswordBase64,
		Uuid:           uuid,
	}
	_, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	return nil, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*repository.User, error) {
	var user repository.User
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return &user, nil
}
func (s *UserStore) ByID(ctx context.Context, id uint) (*repository.User, error) {
	var user repository.User
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return &user, nil
}
