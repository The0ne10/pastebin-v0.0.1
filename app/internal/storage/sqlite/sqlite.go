package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Storage struct {
	db *sql.DB
}

// NewConnect in database
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем соединение с базой данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(alias string, uid uuid.UUID, timeTTL string) error {
	const op = "storage.sqlite.Save"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO aliases (alias, uuid, time_ttl) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, alias, uid, timeTTL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) FindUUID(alias string) (string, error) {
	const op = "storage.sqlite.Find"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var uid string

	err := s.db.QueryRowContext(ctx, "SELECT uuid FROM aliases WHERE alias = ?", alias).Scan(&uid)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}
