package bot

import (
	"context"
	events "telegram-informer/internal/bot/handlers"
	addeventtext "telegram-informer/internal/bot/handlers/add-event-text"
	createevent "telegram-informer/internal/bot/handlers/create-event"
	deletealleventstoday "telegram-informer/internal/bot/handlers/delete-all-events-today"
	deleteeventbyid "telegram-informer/internal/bot/handlers/delete-event-by-id"
	geteventbyid "telegram-informer/internal/bot/handlers/get-event-by-id"
	geteventtoday "telegram-informer/internal/bot/handlers/get-event-today"
	mainmenu "telegram-informer/internal/bot/handlers/main_menu"
	"telegram-informer/internal/bot/state"
	"telegram-informer/internal/domain"
	"time"

	"github.com/go-telegram/bot"
)

const empty = ""

type StorageService interface {
	AddEvent(ctx context.Context, userId int, title string, time time.Time, timeToNotify time.Time) error
	DeleteEvent(ctx context.Context, userId int, eventId int) error
	DeleteEventFromToday(ctx context.Context, userId int) error
	GetEvents(ctx context.Context, userId int) ([]domain.Event, error)
	GetEvent(ctx context.Context, userId int, id int) (domain.Event, error)
	GetEventsTodayFromUser(ctx context.Context, userId int) ([]domain.Event, error)
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
	addEventText := addeventtext.NewHandle(storage, stateStore)
	mainMenu := mainmenu.NewHandle()
	createEvent := createevent.NewHandle(stateStore)
	deleteEventById := deleteeventbyid.NewHandle(storage)
	getEventToday := geteventtoday.NewHandle(storage)
	deleteAllEventsToday := deletealleventstoday.NewHandle(storage)

	b.RegisterHandler(bot.HandlerTypeMessageText, events.Start, bot.MatchTypePrefix, mainMenu.Handler)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBAddEvent, bot.MatchTypeExact, createEvent.Handler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBTodayEvents, bot.MatchTypeExact, getEventToday.Handler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBCancelAllTodayEvents, bot.MatchTypeExact, deleteAllEventsToday.Handler)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBGetById, bot.MatchTypePrefix, getEventById.Handle)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBDeleteById, bot.MatchTypePrefix, deleteEventById.Handler)

	b.RegisterHandler(bot.HandlerTypeMessageText, empty, bot.MatchTypePrefix, addEventText.Handle)
}
