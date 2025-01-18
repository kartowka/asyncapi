package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SigninResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r SigninRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
func (s *Server) signinHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[SigninRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}
		user, err := s.store.Users.ByEmail(r.Context(), req.Email)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		if err := user.ComparePassword(req.Password); err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}
		tokenPair, err := s.JWTManager.GenerateTokenPair(user.ID)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		_, err = s.store.RefreshTokens.DeleteUserTokens(r.Context(), user.ID)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		_, err = s.store.RefreshTokens.Create(r.Context(), user.ID, tokenPair.RefreshToken)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		if err := encode(ApiResponse[SigninResponse]{Data: &SigninResponse{
			AccessToken:  tokenPair.AccessToken.Raw,
			RefreshToken: tokenPair.RefreshToken.Raw,
		},
		}, http.StatusOK, w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		return nil
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh_token is required")
	}
	return nil
}
func (s *Server) refreshTokenHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		req, err := decode[RefreshTokenRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}
		currentRefreshToken, err := s.JWTManager.Parse(req.RefreshToken)
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}
		userIDstr, err := currentRefreshToken.Claims.GetSubject()
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}
		userID, err := strconv.ParseUint(userIDstr, 10, 64)
		if err != nil {
			return NewErrWithStatus(http.StatusUnauthorized, err)
		}
		currentRefreshTokenRecord, err := s.store.RefreshTokens.GetByToken(r.Context(), uint(userID), currentRefreshToken)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, sql.ErrNoRows) {
				status = http.StatusUnauthorized
			}
			return NewErrWithStatus(status, err)
		}
		if currentRefreshTokenRecord.ExpiresAt.Before(time.Now()) {
			return NewErrWithStatus(http.StatusUnauthorized, errors.New("refresh token has expired"))
		}
		TokenPair, err := s.JWTManager.GenerateTokenPair(uint(userID))
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		if _, err := s.store.RefreshTokens.DeleteUserTokens(r.Context(), uint(userID)); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		if _, err := s.store.RefreshTokens.Create(r.Context(), uint(userID), TokenPair.RefreshToken); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		if err := encode(ApiResponse[SigninResponse]{Data: &SigninResponse{
			AccessToken:  TokenPair.AccessToken.Raw,
			RefreshToken: TokenPair.RefreshToken.Raw,
		},
		}, http.StatusOK, w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}
		return nil
	})
}
