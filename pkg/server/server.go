package server

import (
	"context"
	"fmt"
	"iFall/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	App    *fiber.App
	Config config.ServerConfig
}

func NewServer(scfg config.ServerConfig, acfg config.AppConfig) *Server {
	app := fiber.New(fiber.Config{
		AppName:      acfg.Name,
		ErrorHandler: ErrorHandler,
	})
	app.Use(
		cors.New(cors.ConfigDefault),
		RequestTimeoutMiddleware(scfg.RequestTimeout),
	)
	server := &Server{
		App:    app,
		Config: scfg,
	}
	return server
}

func (s *Server) Start() {
	if err := s.App.Listen(s.Config.Host + ":" + s.Config.Port); err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}
}

func (s *Server) MustClose() {
	ctx, cancel := context.WithTimeout(context.Background(), s.Config.CloseTimeout)
	defer cancel()
	if err := s.App.ShutdownWithContext(ctx); err != nil {
		panic(fmt.Errorf("failed to close server: %w", err))
	}
}
