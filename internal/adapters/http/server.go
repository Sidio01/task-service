package http

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team26/task/internal/ports"
	httpMiddleware "gitlab.com/g6834/team26/task/pkg/middleware"
)

type Server struct {
	task     ports.Task
	server   *http.Server
	logger   *zerolog.Logger
	listener net.Listener
	port     int
}

func New(l *zerolog.Logger, task ports.Task) (*Server, error) {
	var (
		err error
		s   Server
	)
	s.listener, err = net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal("Failed listen port", err)
	}
	s.task = task
	s.logger = l
	s.port = s.listener.Addr().(*net.TCPAddr).Port

	s.server = &http.Server{
		Handler: s.routes(),
	}

	return &s, nil
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Start() error {
	if err := s.server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(httpMiddleware.LoggerMiddleware(s.logger))
	r.Use(httpMiddleware.RecovererMiddleware(s.logger))
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	// r.Use(middleware.Recoverer)
	// r.Use(middleware.Timeout(60 * time.Second))
	// r.Get("/healthz", s.healthzHandler)
	r.Mount("/task/v1", s.authHandlers())
	return r
}

// func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// }
