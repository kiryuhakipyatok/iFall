package apierr

import (
	"errors"
	"fmt"
	"iFall/pkg/errs"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrBadRequest     = errors.New("bad request")
	ErrInternalServer = errors.New("internal server error")
	ErrRequestTimeout = errors.New("request timeout")
	ErrToManyRequests = errors.New("to many requests")
	ErrInvalidJSON    = errors.New("invalid json")
)

type ApiErr struct {
	Code    int
	Message any
}

func (ae ApiErr) Error() string {
	return fmt.Sprintf("error: %s, code: %d", ae.Message, ae.Code)
}

func NewApiError(code int, err error) ApiErr {
	return ApiErr{
		Code:    code,
		Message: err.Error(),
	}
}

func ValidationError(errors map[string]string) ApiErr {
	return ApiErr{
		Code:    fiber.StatusUnprocessableEntity,
		Message: errors,
	}
}

func ToApiError(err error) ApiErr {
	switch {
	case errors.Is(err, errs.ErrAlreadyExistsBase):
		return AlreadyExists()
	case errors.Is(err, errs.ErrNotFoundBase):
		return NotFound()
	default:
		return InternalServerError()
	}
}

func InvalidJSON() ApiErr {
	return NewApiError(fiber.StatusUnprocessableEntity, ErrInvalidJSON)
}

func InternalServerError() ApiErr {
	return NewApiError(fiber.StatusInternalServerError, ErrInternalServer)
}

func InvalidRequest() ApiErr {
	return NewApiError(fiber.StatusBadRequest, ErrBadRequest)
}

func NotFound() ApiErr {
	return NewApiError(fiber.StatusNotFound, ErrNotFound)
}

func AlreadyExists() ApiErr {
	return NewApiError(fiber.StatusConflict, ErrAlreadyExists)
}

func RequestTimeout() ApiErr {
	return NewApiError(fiber.StatusRequestTimeout, ErrRequestTimeout)
}

func TooManyRequests() ApiErr {
	return NewApiError(fiber.StatusTooManyRequests, ErrToManyRequests)
}
