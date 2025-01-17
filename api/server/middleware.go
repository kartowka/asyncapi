package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"strconv"

	"github.com/antfley/asyncapi/internal/repository"
	"github.com/antfley/asyncapi/store"
)

func NewLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("http request", "path", r.Method+" "+r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

type userCtxKey struct{}

func ContextWithUser(ctx context.Context, user *repository.User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}
func NewAuthMiddleware(jm JWTManager, userStore *store.UserStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth") {
				next.ServeHTTP(w, r)
				return
			}
			authHeaders := r.Header.Get("Authorization")
			var token string
			if parts := strings.Split(authHeaders, "Bearer "); len(parts) == 2 {
				token = parts[1]
			}
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			parsedToken, err := jm.Parse(token)
			if err != nil {
				slog.Error("error parsing token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !jm.IsAccessToken(parsedToken) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not access token"))
				return
			}
			userIdstr, err := parsedToken.Claims.GetSubject()
			if err != nil {
				slog.Error("failed to extract user id from token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			userID, err := strconv.ParseUint(userIdstr, 10, 32)
			if err != nil {
				slog.Error("failed to convert user id to uint", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := userStore.ByID(r.Context(), uint(userID))
			if err != nil {
				slog.Error("failed to get user by id", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), user)))
		},
		)
	}
}
