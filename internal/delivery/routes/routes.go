package routes

import (
	"iFall/internal/delivery/handlers"

	"github.com/gofiber/fiber/v2"
)

type RoutesSetup struct {
	App         *fiber.App
	UserHandler *handlers.UsersHandler
}

func NewRoutesSetup(a *fiber.App, uh *handlers.UsersHandler) *RoutesSetup {
	return &RoutesSetup{
		App:         a,
		UserHandler: uh,
	}
}

func (rs *RoutesSetup) SetupRoutes() {
	rs.UsersRoutes()
}

func (rs *RoutesSetup) UsersRoutes() {
	rs.App.Post("/api/v1/users", rs.UserHandler.CreateUser)
}
