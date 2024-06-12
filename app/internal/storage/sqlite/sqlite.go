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

func (s *Storage) DeleteExpired() ([]string, error) {
	const op = "storage.sqlite.DeleteExpired"

	rows, err := s.db.Query("SELECT uuid FROM aliases WHERE time_ttl < ?", time.Now())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var uids []string
	for rows.Next() {
		var uid string
		if err = rows.Scan(&uid); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		uids = append(uids, uid)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Удаляем записи, которые были возвращены
	if len(uids) > 0 {
		_, err = s.db.Exec("DELETE FROM aliases WHERE time_ttl < ?", time.Now())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return uids, nil
}
