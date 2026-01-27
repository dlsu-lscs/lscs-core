package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
)

// Setup configures the global zerolog logger based on the application config.
// It should be called once at application startup after config.Load().
func Setup(cfg *config.Config) {
	// set log level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// configure output format based on environment
	if cfg.IsDevelopment() {
		// human-readable console output for development
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production (better for log aggregation)
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// Logger returns the global logger for package-level logging.
// Prefer using log.Info(), log.Error(), etc. directly from zerolog for simpler code.
func Logger() *zerolog.Logger {
	return &log.Logger
}
