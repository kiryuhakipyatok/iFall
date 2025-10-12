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

func TestUserRepository_DropChatId(t *testing.T) {
	tests := []struct {
		testName       string
		telegram       string
		chatId         int64
		expectedResult error
	}{
		{
			testName:       "success dropping",
			telegram:       "kirtg",
			chatId:         123123,
			expectedResult: nil,
		},
		{
			testName:       "not found tg",
			telegram:       "tg",
			chatId:         123123,
			expectedResult: errs.ErrNotFoundBase,
		},
		{
			testName:       "not found chatid",
			telegram:       "kirtg",
			chatId:         111111,
			expectedResult: errs.ErrNotFoundBase,
		},
		{
			testName:       "not found tg and chatid",
			telegram:       "tg",
			chatId:         111111,
			expectedResult: errs.ErrNotFoundBase,
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

	if _, err := storage.DB.Exec("INSERT INTO users (id, name, email, telegram, chat_id) VALUES($1, $2, $3, $4, $5)", uuid.New(), "kir", "kir@gmail.com", "kirtg", 123123); err != nil {
		t.Fatalf("failed to insert test user data in the table: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			repo := NewUserRepository(storage)
			err := repo.DropChatId(context.Background(), tt.telegram, tt.chatId)
			if tt.expectedResult == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedResult)
			}
		})
	}
}

func TestUserRepository_SetChatId(t *testing.T) {
	tests := []struct {
		testName       string
		telegram       string
		chatId         int64
		expectedResult error
	}{
		{
			testName:       "success setting",
			telegram:       "kirtg",
			chatId:         456456,
			expectedResult: nil,
		},
		{
			testName:       "not found tg",
			telegram:       "tg",
			chatId:         111111,
			expectedResult: errs.ErrNotFoundBase,
		},
		{
			testName:       "already exists",
			telegram:       "santg",
			chatId:         123123,
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

	if _, err := storage.DB.Exec("INSERT INTO users (id, name, email, telegram, chat_id) VALUES($1, $2, $3, $4, $5)", uuid.New(), "sanya", "san@gmail.com", "santg", 123123); err != nil {
		t.Fatalf("failed to insert test user data in the table: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			repo := NewUserRepository(storage)
			err := repo.SetChatId(context.Background(), tt.telegram, tt.chatId)
			if tt.expectedResult == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedResult)
			}
		})
	}
}

func TestUserRepository_FetchContacts(t *testing.T) {
	tests := []struct {
		testName       string
		expectedResult []models.Contacts
		expectedError  error
	}{
		{
			testName: "success fetching",
			expectedResult: []models.Contacts{
				{
					Email:  "kiremail",
					ChatId: utils.Int64ToPtr(123123),
				},
				{
					Email: "gusemail",
				},
			},
			expectedError: nil,
		},
		{
			testName:       "success empty fetching",
			expectedResult: []models.Contacts{},
			expectedError:  nil,
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

	if _, err := storage.DB.Exec("INSERT INTO users (id, name, email, telegram, chat_id) VALUES($1, $2, $3, $4, $5)", uuid.New(), "kir", "kiremail", "kirtg", 123123); err != nil {
		t.Fatalf("failed to insert test user data in the table: %v", err)
	}

	if _, err := storage.DB.Exec("INSERT INTO users (id, name, email, telegram, chat_id) VALUES($1, $2, $3, $4, $5)", uuid.New(), "sanya", "gusemail", "gustg", nil); err != nil {
		t.Fatalf("failed to insert test user data in the table: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			defer func() {
				if _, err := storage.DB.Exec("DELETE FROM users"); err != nil {
					t.Fatalf("failed to delete users: %v", err)
				}
			}()
			repo := NewUserRepository(storage)
			contacts, err := repo.FetchContacts(context.Background())
			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.ElementsMatch(t, contacts, tt.expectedResult)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}

func TestUserRepository_CheckChatId(t *testing.T) {
	tests := []struct {
		testName       string
		telegram       string
		chatId         int64
		expectedResult bool
		expectedError  error
	}{
		{
			testName:       "success exist",
			telegram:       "tg1",
			chatId:         123123,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			testName:       "success not exist",
			telegram:       "tg2",
			chatId:         456456,
			expectedResult: true,
			expectedError:  nil,
		},
		{
			testName:       "tg not found",
			telegram:       "tg3",
			chatId:         456456,
			expectedResult: false,
			expectedError:  errs.ErrNotFoundBase,
		},
		{
			testName:       "incorrect chatid",
			telegram:       "tg1",
			chatId:         456456,
			expectedResult: false,
			expectedError:  errIncorrectChatId,
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

	if _, err := storage.DB.Exec("INSERT INTO users (id, name, email, telegram, chat_id) VALUES($1, $2, $3, $4, $5)", uuid.New(), "kir", "kiremail", "tg1", 123123); err != nil {
		t.Fatalf("failed to insert test user data in the table: %v", err)
	}

	if _, err := storage.DB.Exec("INSERT INTO users (id, name, email, telegram, chat_id) VALUES($1, $2, $3, $4, $5)", uuid.New(), "sanya", "gusemail", "tg2", nil); err != nil {
		t.Fatalf("failed to insert test user data in the table: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			repo := NewUserRepository(storage)
			exist, err := repo.CheckChatId(context.Background(), "test-op", tt.telegram, tt.chatId)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}
			assert.Equal(t, tt.expectedResult, exist)
		})
	}
}
