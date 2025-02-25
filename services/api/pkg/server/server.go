package server

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/Autherain/saltyChat/internal/health"
	"github.com/Autherain/saltyChat/internal/utils/logger"
	"github.com/Autherain/saltyChat/pkg/store"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/restest"
	"github.com/nats-io/nats.go"
)

type Server struct {
	Service         *res.Service
	log             *logger.Logger
	healthChecker   *health.HealthChecker
	wg              sync.WaitGroup
	shutdownTimeout time.Duration

	store *store.Store
}

type Option func(*Server)

// New creates a new server instance with the given options
func New(options ...Option) *Server {
	s := &Server{}

	for _, option := range options {
		option(s)
	}

	if s.Service == nil {
		panic("server requires a RES service")
	}

	s.addRESHandlers()

	if s.log == nil {
		s.log = logger.NewDefault()
	}

	return s
}

// WithService sets the RES service
func WithService(service *res.Service) Option {
	return func(s *Server) {
		s.Service = service
	}
}

// WithLogger sets the logger
func WithLogger(logger *logger.Logger) Option {
	return func(s *Server) {
		s.log = logger
	}
}

// WithHealthChecker sets the health checker
func WithHealthChecker(healthChecker *health.HealthChecker) Option {
	return func(s *Server) {
		s.healthChecker = healthChecker
	}
}

// WithShutdownTimeout sets the shutdown timeout
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func WithStore(store *store.Store) Option {
	return func(s *Server) {
		s.store = store
	}
}

func (s *Server) Start(ctx context.Context, natsConn *nats.Conn) error {
	s.log.Info("Starting application")

	errChan := make(chan error, 1)

	// Start health checker if enabled
	if s.healthChecker != nil {
		s.startHealthChecker(errChan)
	}

	// Start service
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.Service.Serve(natsConn); err != nil {
			s.log.Error("Service error", slog.Any("error", err))
			errChan <- err
		}
	}()

	return s.handleShutdown(ctx, errChan)
}

func (s *Server) startHealthChecker(errChan chan error) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.healthChecker.Start(); err != nil {
			s.log.Error("Health checker error", slog.Any("error", err))
			errChan <- err
		}
	}()
}

func (s *Server) handleShutdown(ctx context.Context, errChan chan error) error {
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		s.log.Info("Context cancelled")
		return s.shutdown()
	}
}

func (s *Server) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if s.healthChecker != nil {
		s.healthChecker.Stop()
	}

	if err := s.Service.Shutdown(); err != nil {
		s.log.Error("Error stopping RES service", slog.Any("error", err))
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.log.Info("Service stopped gracefully")
		return nil
	case <-shutdownCtx.Done():
		s.log.Error("Service shutdown timed out")
		return shutdownCtx.Err()
	}
}

func (s *Server) addRESHandlers() {
	s.registerRoomRoutes()
}

func newTestSession(t *testing.T, service *res.Service) *restest.Session {
	t.Helper()
	return restest.NewSession(t, service, func(c *restest.SessionConfig) { c.TimeoutDuration = 10 * time.Second })
}
