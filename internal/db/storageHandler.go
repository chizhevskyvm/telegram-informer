package db

import (
	"context"
	"telegram-informer/internal/domain"
	"time"
)

type StorageHandler interface {
	AddEvent(ctx context.Context, userId int, title string, time time.Time, timeToNotify time.Time) error
	DeleteEvent(ctx context.Context, userId int, eventId string) error
	GetEvents(ctx context.Context, userId int) ([]domain.Event, error)
}
