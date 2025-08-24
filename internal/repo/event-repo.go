package repo

import (
	"context"
	"database/sql"
	"fmt"
	"telegram-informer/internal/domain"
	"time"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) AddEvent(ctx context.Context, userId int, title string, t time.Time, timeToNotify time.Time) error {
	const q = `INSERT INTO events (user_id, title, time, timetonotify) VALUES ($1, $2, $3, $4)`
	if _, err := r.db.ExecContext(ctx, q, userId, title, t, timeToNotify); err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, userId int, eventId int) error {
	const q = `DELETE FROM events WHERE user_id = $1 AND id = $2`
	if _, err := r.db.ExecContext(ctx, q, userId, eventId); err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	return nil
}

func (r *EventRepository) GetEvent(ctx context.Context, userId int, id int) (domain.Event, error) {
	const base = selectClause + ` WHERE user_id = $1 AND id = $2`
	row := r.db.QueryRowContext(ctx, base, userId, id)

	ev, err := scanEvent(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Event{}, err
		}
		return domain.Event{}, fmt.Errorf("query row: %w", err)
	}
	return ev, nil
}

func (r *EventRepository) GetEventsToday(ctx context.Context) ([]domain.Event, error) {
	const where = `WHERE time::date = CURRENT_DATE AND is_sent = false`
	return r.queryEvents(ctx, where)
}

func (r *EventRepository) GetEventsTodayFromUser(ctx context.Context, userId int) ([]domain.Event, error) {
	const where = `WHERE user_id = $1 AND time::date = CURRENT_DATE`
	return r.queryEventsWithArgs(ctx, where, userId)
}

func (r *EventRepository) GetEvents(ctx context.Context, userId int) ([]domain.Event, error) {
	const where = `WHERE user_id = $1`
	return r.queryEventsWithArgs(ctx, where, userId)
}

func (r *EventRepository) DeleteEventFromToday(ctx context.Context, userId int) error {
	const q = `DELETE FROM events WHERE user_id = $1 AND time::date = CURRENT_DATE`
	if _, err := r.db.ExecContext(ctx, q, userId); err != nil {
		return fmt.Errorf("delete today events: %w", err)
	}
	return nil
}

const selectClause = `SELECT id, user_id, title, time, timetonotify FROM events`

func (r *EventRepository) queryEvents(ctx context.Context, where string) ([]domain.Event, error) {
	rows, err := r.db.QueryContext(ctx, selectClause+" "+where)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	return scanEvents(rows)
}

func (r *EventRepository) queryEventsWithArgs(ctx context.Context, where string, args ...any) ([]domain.Event, error) {
	rows, err := r.db.QueryContext(ctx, selectClause+" "+where, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	return scanEvents(rows)
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanEvent(rs rowScanner) (domain.Event, error) {
	var e domain.Event
	if err := rs.Scan(&e.ID, &e.UserID, &e.Title, &e.Notification, &e.TimeToNotify); err != nil {
		return domain.Event{}, err
	}
	return e, nil
}

func scanEvents(rows *sql.Rows) ([]domain.Event, error) {
	events := make([]domain.Event, 0, 16)
	for rows.Next() {
		ev, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		events = append(events, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return events, nil
}
