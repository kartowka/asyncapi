package repository

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (u *User) ComparePassword(password string) error {
	hashedPassword, err := base64.StdEncoding.DecodeString(u.HashedPassword)
	if err != nil {
		return fmt.Errorf("could not decode hashed password: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}
