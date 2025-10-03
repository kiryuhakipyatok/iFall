package services

import (
	"context"
	"iFall/internal/domain/models"
	"iFall/internal/domain/repositories"
	"iFall/pkg/errs"
	"iFall/pkg/logger"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, name, email, telegram string) error
}

type userService struct {
	UserRepository repositories.UserRepository
	Logger         *logger.Logger
}

func NewUserService(ur repositories.UserRepository, l *logger.Logger) UserService {
	return &userService{
		UserRepository: ur,
		Logger:         l,
	}
}

func (us *userService) Create(ctx context.Context, name, email, telegram string) error {
	op := "userService.Create"
	log := us.Logger.AddOp(op)
	log.Info("creating user")
	user := &models.User{
		Id:       uuid.New(),
		Name:     name,
		Email:    email,
		Telegram: telegram,
	}
	if err := us.UserRepository.Create(ctx, user); err != nil {
		log.Error("failed to create user", logger.Err(err))
		return errs.NewAppError(op, err)
	}
	return nil
}
