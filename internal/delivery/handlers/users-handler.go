package handlers

import (
	"iFall/internal/delivery/apierr"
	"iFall/internal/domain/services"
	"iFall/internal/dto"
	"iFall/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UsersHandler struct {
	UserService services.UserService
	Validator   *validator.Validator
}

func NewUsersHandler(us services.UserService, v *validator.Validator) *UsersHandler {
	return &UsersHandler{
		UserService: us,
		Validator:   v,
	}
}

func (uh *UsersHandler) CreateUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := dto.CreateUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		return apierr.InvalidJSON()
	}
	if err := uh.Validator.Validate.Struct(req); err != nil {
		return apierr.InvalidRequest()
	}
	if err := uh.UserService.Create(ctx, req.Name, req.Email, req.Telegram); err != nil {
		return apierr.ToApiError(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
