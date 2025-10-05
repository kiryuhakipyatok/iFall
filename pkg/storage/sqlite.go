package storage

import (
	"context"
	"database/sql"
	"iFall/internal/config"

	_ "modernc.org/sqlite"

	"fmt"
)

type Storage struct {
	DB *sql.DB
}

func MustConnect(cfg config.StorageConfig) *Storage {
	path := cfg.Path
	if path == "" {
		panic("storage path is empty")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		panic(fmt.Errorf("failed to connect to sqlite: %w", err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
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
