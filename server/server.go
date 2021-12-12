package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Server is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type Server struct {
	address  string
	listener net.Listener
	logger   *zap.Logger
	server   *http.Server
}

type Options struct {
	Host   string
	Logger *zap.Logger
	Router chi.Router
	Port   int
	Server *http.Server
}

// New configures an instance of the application with helpful defaults.
func New(opts Options) *Server {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))

	if opts.Server == nil {
		opts.Server = &http.Server{
			Addr:    address,
			Handler: opts.Router,
		}
	}

	return &Server{
		address: address,
		logger:  opts.Logger,
		server: &http.Server{
			Handler: opts.Router,
		},
	}
}

func (s *Server) Listen(ctx context.Context) error {
	if s.address == "" {
		s.address = ":http"
	}

	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", s.address)
	if err != nil {
		return fmt.Errorf("error creating listener: %w", err)
	}

	s.listener = listener

	return nil
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info("starting", zap.String("address", s.address))
	go func() {
		if err := s.server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server closed unexpectedly", zap.Error(err))
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}

func (s *Server) URL() string {
	if s.listener == nil {
		return ""
	}

	return s.listener.Addr().String()
}
