package addeventtext

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	eventsstate "telegram-informer/internal/bot/event-state"
	"telegram-informer/internal/bot/ui/texts"
	"telegram-informer/internal/infra/cache"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
)

// Константы для хранения ключей и UI
const (
	titleValue = "title_value"
	dateValue  = "date_value"
	timeValue  = "time_value"
	stateTTL   = 10 * time.Minute
)

type Handle struct {
	eventService EventService
	cache        Cache
}

func NewHandle(eventService EventService, cache Cache) *Handle {
	return &Handle{eventService: eventService, cache: cache}
}

type Cache interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type EventService interface {
	AddEvent(ctx context.Context, userId int, title string, time time.Time, timeToNotify time.Time) error
}

func (h *Handle) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userInput := update.Message.Text
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	state, err := h.cache.Get(eventsstate.GetUserStateKey(userID))
	if err != nil && !errors.Is(err, redis.Nil) {
		fmt.Println("Ошибка при получении состояния:", err)
		return
	}

	if eventsstate.IsAddEventState(state, userID) {
		return
	}

	dataKey := eventsstate.GetUserStateDataKey("addEvent", strconv.FormatInt(userID, 10))
	data, err := cache.GetTyped[map[string]string](h.cache, dataKey)
	if err != nil {
		data = map[string]string{}
	}

	switch state {
	case eventsstate.GetUserStageState(eventsstate.StateAddEventTitle, userID):
		h.addEventTitleState(userID, userInput, data, dataKey)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgAskDate})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventDate, userID):
		h.addEventDateState(userID, userInput, data, dataKey)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgAskTime})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventTime, userID):
		h.addEventTimeState(userID, userInput, data, dataKey)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgConfirm})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventDone, userID):
		h.addEventDoneState(ctx, userID, data, dataKey)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgCreated})
	}
}

func (h *Handle) addEventTitleState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[titleValue] = userInput
	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventDate, userID), stateTTL)
	_ = cache.SetTyped(h.cache, dataKey, data, stateTTL)
}

func (h *Handle) addEventDateState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[dateValue] = userInput
	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventTime, userID), stateTTL)
	_ = cache.SetTyped(h.cache, dataKey, data, stateTTL)
}

func (h *Handle) addEventTimeState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[timeValue] = userInput
	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventDone, userID), stateTTL)
	_ = cache.SetTyped(h.cache, dataKey, data, stateTTL)
}

func (h *Handle) addEventDoneState(ctx context.Context, userID int64, data map[string]string, dataKey string) {
	_ = h.cache.Delete(eventsstate.GetUserStateKey(userID))
	_ = h.cache.Delete(dataKey)

	dateParsed, _ := time.Parse("2006-01-02", data[dateValue])
	timeParsed, _ := time.Parse("15:04", data[timeValue])

	_ = h.eventService.AddEvent(ctx, int(userID), data[titleValue], dateParsed, timeParsed)
}
