package repositories

import (
	"context"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/internal/utils"
	"iFall/pkg/errs"
	"iFall/pkg/storage"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		testName       string
		userData       *models.User
		expectedResult error
	}{
		{
			testName: "success creation with tg",
			userData: &models.User{
				Id:   uuid.New(),
				Name: "sanya",
				Contacts: models.Contacts{
					Email:    "gusemail@gmail.com",
					Telegram: utils.StrToPtr("tg"),
					ChatId:   nil,
				},
			},
			expectedResult: nil,
		},
		{
			testName: "success creation without tg",
			userData: &models.User{
				Id:   uuid.New(),
				Name: "sanya",
				Contacts: models.Contacts{
					Email:    "gusemail@gmail.com",
					Telegram: nil,
					ChatId:   nil,
				},
			},
			expectedResult: nil,
		},
		{
			testName: "already exists with name",
			userData: &models.User{
				Id:   uuid.New(),
				Name: "kir",
				Contacts: models.Contacts{
					Email:    "gusemail@gmail.com",
					Telegram: utils.StrToPtr("tg"),
					ChatId:   nil,
				},
			},
			expectedResult: errs.ErrAlreadyExistsBase,
		},
		{
			testName: "already exists with email",
			userData: &models.User{
				Id:   uuid.New(),
				Name: "sanya",
				Contacts: models.Contacts{
					Email:    "kir@gmail.com",
					Telegram: utils.StrToPtr("tg"),
					ChatId:   nil,
				},
			},
			expectedResult: errs.ErrAlreadyExistsBase,
		},
		{
			testName: "already exists with tg",
			userData: &models.User{
				Id:   uuid.New(),
				Name: "sanya",
				Contacts: models.Contacts{
					Email:    "gusemail@gmail.com",
					Telegram: utils.StrToPtr("kirtg"),
					ChatId:   nil,
				},
			},
			expectedResult: errs.ErrAlreadyExistsBase,
		},
	}
	storage := storage.MustConnect(config.StorageConfig{Path: ":memory:", PingTimeout: time.Second})
	schema := `
		CREATE TABLE IF NOT EXISTS users (
    		id TEXT PRIMARY KEY,
    		name TEXT NOT NULL UNIQUE,
    		email TEXT NOT NULL UNIQUE,
    		telegram TEXT UNIQUE,
    		chat_id INTEGER UNIQUE
		);
	`
	if _, err := storage.DB.Exec(schema); err != nil {
		t.Fatalf("failed to create test users table: %v", err)
	}

	if _, err := storage.DB.Exec("INSERT INTO users (id, name, email, telegram, chat_id) VALUES($1, $2, $3, $4, $5)", uuid.New(), "kir", "kir@gmail.com", "kirtg", nil); err != nil {
		t.Fatalf("failed to insert test user data in the table: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			defer func() {
				if _, err := storage.DB.Exec("DELETE FROM users WHERE id = $1", tt.userData.Id); err != nil {
					t.Fatalf("failed to delete users data from table: %v", err)
				}
			}()

			repo := NewUserRepository(storage)
			err := repo.Create(context.Background(), tt.userData)
			if tt.expectedResult == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedResult)
			}

		})
	}
}
