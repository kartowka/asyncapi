package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ErrWithStatus struct {
	Status int
	err    error
}

func (e *ErrWithStatus) Error() string {
	return e.err.Error()
}
func NewErrWithStatus(status int, err error) error {
	return &ErrWithStatus{Status: status, err: err}
}
func handler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			status := http.StatusInternalServerError
			msg := http.StatusText(status)
			if e, ok := err.(*ErrWithStatus); ok {
				status = e.Status
				msg = http.StatusText(e.Status)
				if status == http.StatusBadRequest || status == http.StatusConflict {
					msg = e.err.Error()
				}
			}

			slog.Error("error executing handler", "error", err, "status", status, "message", msg)
			if err := encode(ApiResponse[struct{}]{Message: msg}, status, w); err != nil {
				slog.Error("error encoding response", "error", err)
			}
		}
	}
}
func encode[T any](v T, status int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}
	return nil
}

type Validator interface {
	Validate() error
}

func decode[T Validator](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("failed to decode request: %w", err)
	}
	if err := v.Validate(); err != nil {
		return v, fmt.Errorf("invalid request: %w", err)
	}
	return v, nil
}
