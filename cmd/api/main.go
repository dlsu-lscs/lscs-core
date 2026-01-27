package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
	"github.com/dlsu-lscs/lscs-core-api/internal/logging"
	"github.com/dlsu-lscs/lscs-core-api/internal/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()
	stop() // Allow Ctrl+C to force shutdown

	log.Info().Msg("shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		return
	}

	logging.Setup(cfg)

	srv := server.NewServer(cfg)

	done := make(chan bool, 1)

	go gracefulShutdown(srv, done)

	log.Info().
		Int("port", cfg.Port).
		Str("env", cfg.GoEnv).
		Msg("starting server")

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("http server error")
	}

	// Wait for the graceful shutdown to complete
	<-done
}
