package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "mysql"),
	}
}

type User struct {
	Id                   int       `db:"id"`
	Email                string    `db:"email"`
	HashedPasswordBase64 string    `db:"hashed_password"`
	CreatedAt            time.Time `db:"created_at"`
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
	const dml = `INSERT INTO users (email, hashed_password) VALUES (?, ?)`

	// Hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)

	// Execute the INSERT statement
	_, err = s.db.ExecContext(ctx, dml, email, hashedPasswordBase64)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	return nil, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	const query = `SELECT * FROM users WHERE email = ?`
	var user User
	if err := s.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return &user, nil
}
func (s *UserStore) ByID(ctx context.Context, userID int) (*User, error) {
	const query = `SELECT * FROM users WHERE id = ?`
	var user User
	if err := s.db.GetContext(ctx, &user, query, userID); err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return &user, nil
}
