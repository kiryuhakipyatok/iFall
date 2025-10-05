package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"fmt"
)

type Storage struct {
	DB *sql.DB
}

func MustConnect(path string) *Storage {
	if path == "" {
		panic("storage path is empty")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		panic(fmt.Errorf("failed to connect to sqlite: %w", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping sqlite: %w", err))
	}
	db.Exec("PRAGMA journal_mode=WAL")
	return &Storage{
		DB: db,
	}
}

func (s *Storage) MustClose() {
	if err := s.DB.Close(); err != nil {
		panic(fmt.Errorf("failed to close sqlite: %w", err))
	}
}
