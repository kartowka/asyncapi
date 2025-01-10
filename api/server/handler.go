package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

type ApiResponse[T any] struct {
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignupRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
func (s *Server) signupHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[SignupRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}
		existingUser, err := s.store.Users.ByEmail(r.Context(), req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		if existingUser != nil {
			return NewErrWithStatus(http.StatusConflict, fmt.Errorf("user with email %s already exists", req.Email))
		}
		_, err = s.store.Users.CreateUser(r.Context(), req.Email, req.Password)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		encode(ApiResponse[struct{}]{Message: "user created"}, http.StatusCreated, w)
		return nil
	})
}
