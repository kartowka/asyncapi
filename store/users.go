package store

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

type User struct {
	ID                   uint      `gorm:"primaryKey"`
	UUID                 uuid.UUID `gorm:"type:char(36);primary_key"`
	Email                string    `gorm:"column:email;unique"`
	HashedPasswordBase64 string    `gorm:"column:hashed_password"`
	CreatedAt            time.Time `gorm:"column:created_at"`
}

func (u *User) ComparePassword(password string) error {
	hashedPassword, err := base64.StdEncoding.DecodeString(u.HashedPasswordBase64)
	if err != nil {
		return fmt.Errorf("could not decode hashed password: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)
	uuid := uuid.New()
	user := User{Email: email, HashedPasswordBase64: hashedPasswordBase64, UUID: uuid}
	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	return nil, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return &user, nil
}
func (s *UserStore) ByID(ctx context.Context, userID int) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return &user, nil
}
