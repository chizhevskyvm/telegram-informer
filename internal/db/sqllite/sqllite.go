package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"telegram-informer/internal/domain"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return &Storage{db}, nil
}

func (s *Storage) AddEvent(ctx context.Context, userId string, time time.Time, information string) {
	stmt, err := s.db.Prepare("INSERT INTO events SET user_id = ?, time = ?, information = ?")
	if err != nil {
		fmt.Printf("scan: %w", err)
	}

	row := stmt.QueryRowContext(ctx, userId, time, information)
	err = row.Scan()
	if err != nil {
		fmt.Printf("scan: %w", err)
	}
}

func (s *Storage) DeleteEvent(ctx context.Context, userId int, eventId string) error {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM events WHERE user_id = ? AND id = ?")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userId, eventId)
	err = row.Scan()
	if err != nil {
		fmt.Printf("scan: %w", err)
	}

	return nil
}

func (s *Storage) GetEvents(ctx context.Context, userId int) ([]domain.Event, error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT id, user_id, title, time FROM events WHERE user_id = ?")
	if err != nil {
		return nil, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userId)
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
