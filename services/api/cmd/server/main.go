package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Autherain/saltyChat/environment"
	"github.com/Autherain/saltyChat/internal/health"
	"github.com/Autherain/saltyChat/internal/utils/logger"
	"github.com/Autherain/saltyChat/pkg/server"
	"github.com/Autherain/saltyChat/pkg/store"
	"github.com/jirenius/go-res"

	_ "github.com/lib/pq"
)

const serviceName = "saltyChat"

func main() {
	// Parse environment variables
	variables := environment.Parse()

	// Initialize logger with the adapter for resgate
	log := logger.NewLogger(logger.Config{
		Format:    variables.LogFormat,
		Level:     variables.LogLevel,
		AddSource: variables.LogSource,
	})
	slog.SetDefault(log.SlogLogger())

	// Initialize NATS connection
	natsConn := environment.MustInitNATSConn(variables)
	defer natsConn.Close()

	// Initialize service first
	service := res.NewService(serviceName)
	// Use the Adapter instead of the logger directly
	service.SetLogger(log.Adapter)
	service.SetInChannelSize(variables.ServiceInChannelSize)
	service.SetWorkerCount(variables.ServiceWorkerCount)

	// Rest of your code remains the same
	healthChecker := health.New(
		natsConn,
		health.NewVersionInfo(variables.Env),
		health.WithNATSCheck(natsConn),
		health.WithInterval(variables.HealthCheckInterval),
		health.WithTimeout(variables.HealthCheckTimeout),
		health.WithSubject(variables.HealthCheckSubject),
		health.WithServiceName(serviceName),
	)

	dbConn := environment.MustInitPGSQLDB(variables)
	store := store.NewStore(store.WithDB(dbConn))

	// Create server with all dependencies
	srv := server.New(
		server.WithService(service),
		server.WithLogger(log), // This is fine as is since server presumably uses the full logger
		server.WithHealthChecker(healthChecker),
		server.WithShutdownTimeout(variables.ShutdownTimeout),
		server.WithStore(store),
	)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Info("Received shutdown signal", "signal", sig)
		cancel()
	}()

	// Start server
	if err := srv.Start(ctx, natsConn); err != nil {
		log.Error("Server error", slog.Any("error", err))
		os.Exit(1)
	}
}
