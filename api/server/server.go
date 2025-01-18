package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/antfley/asyncapi/config"
	"github.com/antfley/asyncapi/store"
)

type Server struct {
	config     *config.Config
	logger     *slog.Logger
	store      *store.Store
	JWTManager *JWTManager
	middleware []func(http.Handler) http.Handler
}

func New(config *config.Config, logger *slog.Logger, store *store.Store, jwtManager *JWTManager) *Server {
	return &Server{
		config:     config,
		logger:     logger,
		store:      store,
		JWTManager: jwtManager,
		middleware: []func(http.Handler) http.Handler{},
	}
}
func (s *Server) Ping() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		encode(ApiResponse[struct{}]{Message: "pong"}, http.StatusOK, w)
		return nil
	})
}
func (s *Server) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.Ping())
	mux.HandleFunc("POST /auth/signup", s.signupHandler())
	mux.HandleFunc("POST /auth/signin", s.signinHandler())
	mux.HandleFunc("POST /auth/refresh", s.refreshTokenHandler())
	return mux
}
func (s *Server) Use(mw func(http.Handler) http.Handler) {
	s.middleware = append(s.middleware, mw)
}
func (s *Server) applyMiddleware(h http.Handler) http.Handler {
	for _, mw := range s.middleware {
		h = mw(h)
	}
	return h
}
func (s *Server) Run(ctx context.Context) error {
	mux := s.router()
	s.Use(NewLoggerMiddleware(s.logger))
	s.Use(NewAuthMiddleware(*s.JWTManager, s.store.Users))
	handler := s.applyMiddleware(mux)
	server := &http.Server{
		Addr:    net.JoinHostPort("", s.config.PORT),
		Handler: handler,
	}
	go func() {
		s.logger.Info("server started", "port", s.config.PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("server error", "error", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("server shutdown error", "error", err)
		}
	}()
	wg.Wait()
	return nil
}
