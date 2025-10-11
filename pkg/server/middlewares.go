package server

import (
	"context"
	"errors"
	"iFall/internal/delivery/apierr"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return c.Status(fe.Code).JSON(fe)
	}
	var apiErr apierr.ApiErr
	if errors.As(err, &apiErr) {
		return c.Status(apiErr.Code).JSON(apiErr)
	}
	internalErr := apierr.InternalServerError()
	return c.Status(internalErr.Code).JSON(internalErr)
}

func RequestTimeoutMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), timeout)
		defer cancel()
		c.SetUserContext(ctx)
		err := c.Next()
		if err == context.DeadlineExceeded {
			return apierr.RequestTimeout()
		}
		return err
	}
}
