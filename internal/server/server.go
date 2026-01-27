package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/dlsu-lscs/lscs-core-api/internal/committee"
	"github.com/dlsu-lscs/lscs-core-api/internal/config"
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/member"
	"github.com/labstack/echo/v4"
)

type Server struct {
	port int
	cfg  *config.Config

	db database.Service

	authHandler      *auth.Handler
	memberHandler    *member.Handler
	committeeHandler *committee.Handler
}

func NewServer(cfg *config.Config) *http.Server {
	dbService := database.New(cfg)

	NewServer := &Server{
		port:             cfg.Port,
		cfg:              cfg,
		db:               dbService,
		authHandler:      auth.NewHandler(auth.NewService(cfg.JWTSecret, cfg), dbService),
		memberHandler:    member.NewHandler(dbService),
		committeeHandler: committee.NewHandler(dbService),
	}

	// Declare Server config
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	NewServer.RegisterRoutes(e)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      e,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
