package repositories

import (
	"context"
	"errors"
	"iFall/internal/domain/models"
	"iFall/pkg/errs"
	"iFall/pkg/storage"
)

//go:generate mockgen -source=users-repo.go -destination=mocks/users-repo-mock.go
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FetchContacts(ctx context.Context) ([]models.Contacts, error)
	SetChatId(ctx context.Context, telegram string, chatId int64) error
	DropChatId(ctx context.Context, telegram string, chatId int64) error
}

type userRepository struct {
	Storage *storage.Storage
}

func NewUserRepository(s *storage.Storage) UserRepository {
	return &userRepository{
		Storage: s,
	}
}

const usersRepo = "userRepository."

func (ur *userRepository) Create(ctx context.Context, user *models.User) error {
	op := usersRepo + "Create"
	query := "INSERT INTO users (id, name, email, telegram) VALUES ($1, $2, $3, $4)"
	if _, err := ur.Storage.DB.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Telegram); err != nil {
		if storage.ErrorAlreadyExists(err) {
			return errs.ErrAlreadyExists(op, err)
		}
		return errs.NewAppError(op, err)
	}
	return nil
}

func (ur *userRepository) DropChatId(ctx context.Context, telegram string, chatId int64) error {
	op := usersRepo + "DropChatId"
	exist, err := ur.checkChatId(ctx, op, telegram, chatId)
	if err != nil {
		return err
	}
	if exist {
		return errs.ErrNotFound(op)
	}
	query := "UPDATE users SET chat_id = null WHERE telegram = $1 AND chat_id = $2"
	res, err := ur.Storage.DB.ExecContext(ctx, query, telegram, chatId)
	if err != nil {
		return errs.NewAppError(op, err)
	}
	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		return errs.ErrNotFound(op)
	}
	return nil
}

func (ur *userRepository) SetChatId(ctx context.Context, telegram string, chatId int64) error {
	op := usersRepo + "SetChatId"
	exist, err := ur.checkChatId(ctx, op, telegram, chatId)
	if err != nil {
		return err
	}
	if !exist {
		return errs.ErrAlreadyExists(op, errors.New("same chat_id already exists"))
	}

	uQuery := "UPDATE users SET chat_id = $1 WHERE telegram = $2"
	if _, err := ur.Storage.DB.ExecContext(ctx, uQuery, chatId, telegram); err != nil {
		return errs.NewAppError(op, err)
	}

	return nil
}

func (ur *userRepository) FetchContacts(ctx context.Context) ([]models.Contacts, error) {
	op := usersRepo + "FetchContacts"
	query := "SELECT email, chat_id FROM users"
	contacts := []models.Contacts{}
	res, err := ur.Storage.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, errs.NewAppError(op, err)
	}
	defer res.Close()
	for res.Next() {
		var contact models.Contacts
		if err := res.Scan(
			&contact.Email,
			&contact.ChatId,
		); err != nil {
			return nil, errs.NewAppError(op, err)
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (ur *userRepository) checkChatId(ctx context.Context, op, telegram string, chatId int64) (bool, error) {
	var cid *int64
	cQuery := "SELECT chat_id FROM users WHERE telegram = $1"
	if err := ur.Storage.DB.QueryRowContext(ctx, cQuery, telegram).Scan(&cid); err != nil {
		if err == storage.ErrNotFound() {
			return false, errs.ErrNotFound(op)
		}
		return false, errs.NewAppError(op, err)
	}

	if cid != nil && *cid == chatId {
		return false, nil
	}

	return true, nil
}
