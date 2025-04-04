package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"telegram-informer/internal/domain"
	"time"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) RunMigrations() error {
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("db driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func New(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return &Storage{db}, nil
}

func (s *Storage) AddEvent(ctx context.Context, userId int, title string, time time.Time, timeToNotify time.Time) error {
	query := `INSERT INTO events (user_id, title, time, timetonotify) VALUES ($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, query, userId, title, time, timeToNotify)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, userId int, eventId string) error {
	query := `DELETE FROM events WHERE user_id = $1 AND id = $2`
	_, err := s.db.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	return nil
}

func (s *Storage) GetEvents(ctx context.Context, userId int) ([]domain.Event, error) {
	query := `SELECT id, user_id, title, time, timetonotify FROM events WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		err := rows.Scan(&e.ID, &e.UserID, &e.Title, &e.Notification, &e.TimeToNotify)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return events, nil
}
