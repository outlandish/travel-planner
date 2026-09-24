package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
	"travelplanner/internal/domain"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func (s *Store) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func CacheKey(holidayParams domain.HolidayParams) string {
	return fmt.Sprintf(
		"%d:%d|%d:%s|%d:%s",
		len(strconv.Itoa(holidayParams.Budget)),
		holidayParams.Budget,
		len(string(holidayParams.HolidayType)),
		string(holidayParams.HolidayType),
		len(string(holidayParams.Nature)),
		string(holidayParams.Nature),
	)
}

func (s *Store) GetResponse(ctx context.Context, key string) (response string, found bool, err error) {
	var res string
	row := s.db.QueryRowContext(ctx, "SELECT response FROM responses WHERE cache_key = ?", key)
	err = row.Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}

		return "", false, err
	}

	return res, true, nil
}

func (s *Store) SaveResponse(ctx context.Context, key string, response string) error {
	createdAt := time.Now().UTC().Format(time.RFC3339)

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO responses (cache_key, response, created_at)
		VALUES (?, ?, ?)
		ON CONFLICT(cache_key)
		DO UPDATE SET
			response = excluded.response,
		    created_at = excluded.created_at
	`, key, response, createdAt)

	return err
}

func NewStore(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("DB path is empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS responses (
    		cache_key TEXT PRIMARY KEY,
    		response  TEXT NOT NULL,
    		created_at TEXT NOT NULL
  		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}
