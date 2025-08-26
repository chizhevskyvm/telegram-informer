package state

import (
	"fmt"
	"telegram-informer/internal/domain"
	"telegram-informer/internal/infra/cache"
	"time"
)

const ttl = 10 * time.Minute

type Cache interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Store struct {
	c Cache
}

func NewStore(c Cache) *Store {
	return &Store{c: c}
}

func (s *Store) GetState(userID int64) (State, error) {
	val, err := s.c.Get(UserStateKey(userID))
	return State(val), err
}

func (s *Store) SetCreateEventState(userID int64) error {
	return s.setState(userID, CreateEventState(userID))
}

func (s *Store) SetEventAddTimeState(userID int64) error {
	return s.setState(userID, AddEventTimeState(userID))
}

func (s *Store) SetEventAddDateState(userID int64) error {
	return s.setState(userID, AddEventDateState(userID))
}

func (s *Store) SetEventAddTitleState(userID int64) error {
	return s.setState(userID, AddEventTitleState(userID))
}

func (s *Store) SetDoneState(userID int64) error {
	return s.setState(userID, AddEventDoneState(userID))
}

func (s *Store) setState(userID int64, state State) error {
	return s.c.Set(UserStateKey(userID), string(state), ttl)
}

func (s *Store) ClearState(userID int64) error {
	return s.c.Delete(UserStateKey(userID))
}

func (s *Store) ClearEventData(userID int64) error {
	return s.c.Delete(DataKey(AddEvent, fmt.Sprint(userID)))
}

func (s *Store) GetAddEventData(userID int64) (*domain.EventData, error) {
	key := DataKey(AddEvent, fmt.Sprint(userID))
	data, err := cache.GetTyped[map[string]string](s.c, key)
	if err != nil || data == nil {
		return domain.NewEventData(nil), err
	}
	return domain.NewEventData(data), nil
}

func (s *Store) SetAddEventData(userID int64, data *domain.EventData) error {
	key := DataKey(AddEvent, fmt.Sprint(userID))
	return cache.SetTyped(s.c, key, data.Raw(), ttl)
}
