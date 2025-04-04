package db

import (
	"context"
	"telegram-informer/internal/domain"
	"time"
)

type StorageHandler interface {
	AddEvent(ctx context.Context, userId string, time time.Time, information string)
	DeleteEvent(ctx context.Context, userId int, eventId string) error
	GetEvents(ctx context.Context, userId int) ([]domain.Event, error)
}
