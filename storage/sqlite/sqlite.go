package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhayt/read-adviser-bot/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect database: %w", err)
	}

	return &Storage{db: db}, nil
}

// Save insert page in storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	stmt := "INSERT INTO page(url, user_name) VALUES (?, ?)"

	_, err := s.db.ExecContext(ctx, stmt, p.URL, p.UserName)
	if err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// PickRandom return random page from storage if it is exists.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	stmt := "SELECT url FROM page WHERE user_name = ? ORDER BY RANDOM() LIMIT 1"

	var url string

	err := s.db.QueryRowContext(ctx, stmt, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{URL: url, UserName: userName}, nil
}

// Remove delete page from storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	stmt := "DELETE FROM page WHERE url = ? and user_name = ?"

	_, err := s.db.ExecContext(ctx, stmt, p.URL, p.UserName)
	if err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

// IsExists check is page exist in storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	stmt := "SELECT count(*) FROM page WHERE url = ? AND user_name = ?"

	var count int

	err := s.db.QueryRowContext(ctx, stmt, p.URL, p.UserName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("can't check is page exist: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	stmt := `CREATE TABLE IF NOT EXISTS page(url Text, user_name Text)`

	_, err := s.db.ExecContext(ctx, stmt)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
