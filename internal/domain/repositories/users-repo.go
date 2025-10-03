package repositories

import (
	"context"
	"iFall/internal/domain/models"
	"iFall/pkg/errs"
	"iFall/pkg/storage"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FetchEmails(ctx context.Context) ([]string, error)
}

type userRepository struct {
	Storage *storage.Storage
}

func NewUserRepository(s *storage.Storage) UserRepository {
	return &userRepository{
		Storage: s,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *models.User) error {
	op := "userRepository.Create"
	query := "INSERT INTO users (id, name, email, telegram) VALUES ($1, $2, $3, $4)"
	res, err := ur.Storage.Pool.Exec(ctx, query, user.Id, user.Name, user.Email, user.Telegram)
	if err != nil {
		if storage.ErrorAlreadyExists(err) {
			return errs.ErrAlreadyExists(op, err)
		}
		return errs.NewAppError(op, err)
	}
	if res.RowsAffected() == 0 {
		return errs.ErrNotFound(op)
	}
	return nil
}

func (ur *userRepository) FetchEmails(ctx context.Context) ([]string, error) {
	op := "userRepository.Create"
	query := "SELECT email FROM users"
	emails := []string{}
	res, err := ur.Storage.Pool.Query(ctx, query)
	if err != nil {
		return nil, errs.NewAppError(op, err)
	}
	defer res.Close()
	for res.Next() {
		var email string
		if err := res.Scan(&email); err != nil {
			return nil, errs.NewAppError(op, err)
		}
		emails = append(emails, email)
	}
	return emails, nil
}
