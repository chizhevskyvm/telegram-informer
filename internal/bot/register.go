package bot

import (
	"context"
	events "telegram-informer/internal/bot/handlers"
	addevent "telegram-informer/internal/bot/handlers/add_event"
	deletealleventstoday "telegram-informer/internal/bot/handlers/delete_all_events_today"
	deleteeventbyid "telegram-informer/internal/bot/handlers/delete_event_by_id"
	geteventbyid "telegram-informer/internal/bot/handlers/get_event_by_id"
	geteventtoday "telegram-informer/internal/bot/handlers/get_event_today"
	geteventsactual "telegram-informer/internal/bot/handlers/get_events_actual"
	mainmenu "telegram-informer/internal/bot/handlers/main_menu"
	setcreateeventstate "telegram-informer/internal/bot/handlers/set_create_event_state"
	"telegram-informer/internal/bot/state"
	"telegram-informer/internal/domain"
	"time"

	"github.com/go-telegram/bot"
)

type StorageService interface {
	AddEvent(ctx context.Context, userId int, title string, time time.Time, timeToNotify time.Time) error
	DeleteEvent(ctx context.Context, userId int, eventId int) error
	DeleteEventFromToday(ctx context.Context, userId int) error
	GetEvents(ctx context.Context, userId int) ([]domain.Event, error)
	GetEvent(ctx context.Context, userId int, id int) (domain.Event, error)
	GetEventsTodayFromUser(ctx context.Context, userId int) ([]domain.Event, error)
	GetEventsActual(ctx context.Context, userId int) ([]domain.Event, error)
	GetEventsToday(ctx context.Context) ([]domain.Event, error)
}

type Cache interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

func RegisterHandlers(b *bot.Bot, storage StorageService, cache Cache) {
	stateStore := state.NewStore(cache)

	getEventById := geteventbyid.NewHandle(storage)
	addEvent := addevent.NewHandle(storage, stateStore)
	mainMenu := mainmenu.NewHandle()
	createEvent := setcreateeventstate.NewHandle(stateStore)
	deleteEventById := deleteeventbyid.NewHandle(storage)
	getEventToday := geteventtoday.NewHandle(storage)
	getEventsActual := geteventsactual.NewHandle(storage)
	deleteAllEventsToday := deletealleventstoday.NewHandle(storage)

	b.RegisterHandler(bot.HandlerTypeMessageText, events.MTPStart, bot.MatchTypePrefix, mainMenu.Handler)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBSetCreateEventState, bot.MatchTypeExact, createEvent.Handler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBGetEventToday, bot.MatchTypeExact, getEventToday.Handler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBDeleteAllEventsToday, bot.MatchTypeExact, deleteAllEventsToday.Handler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBGetEventsActual, bot.MatchTypeExact, getEventsActual.Handler)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBGetEventById, bot.MatchTypePrefix, getEventById.Handle)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBDeleteEventById, bot.MatchTypePrefix, deleteEventById.Handler)

	b.RegisterHandler(bot.HandlerTypeMessageText, events.MTCAddEvent, bot.MatchTypeContains, addEvent.Handle)
}
