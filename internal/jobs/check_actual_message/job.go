package worker

import (
	"context"
	"telegram-informer/internal/domain"
)

type PubSub interface {
	Publish(ctx context.Context, subject string, msg []byte) error
}

type EventService interface {
	GetEventsToday(ctx context.Context) ([]domain.Event, error)
}

type Job struct {
	Storage EventService
	PubSub  PubSub
}

func NewJob(storageHandler EventService, pubSub PubSub) *Job {
	return &Job{
		Storage: storageHandler,
		PubSub:  pubSub,
	}
}
