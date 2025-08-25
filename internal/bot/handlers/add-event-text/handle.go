package addeventtext

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	eventsstate "telegram-informer/internal/bot/event-state"
	"telegram-informer/internal/bot/ui/texts"
	"telegram-informer/internal/infra/cache"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
)

const (
	cacheKeyTitleValue = "title_value"
	cacheKeyDateValue  = "date_value"
	cacheKeyTimeValue  = "time_value"

	addEvent = "addEvent"

	stateTTL = 10 * time.Minute
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

	state, err := h.getUserState(userID)
	if err != nil {
		fmt.Println("Ошибка при получении состояния:", err)
		return
	}

	if !eventsstate.IsAddEventState(state, userID) {
		return
	}

	dataKey, data := h.getUserData(userID, addEvent)

	switch state {
	case eventsstate.GetUserStageState(eventsstate.StateAddEventTitle, userID):
		title := strings.TrimSpace(userInput)
		if title == "" {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.ErrTitleEmpty})
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgAskTitle})
			return
		}

		h.addEventTitleState(userID, userInput, data, dataKey)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgAskDate})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventDate, userID):
		dateStr := strings.TrimSpace(userInput)
		if _, err := parseDateLocal(dateStr); err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.ErrDateFormat})
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgAskDate})
			return
		}
		h.addEventDateState(userID, dateStr, data, dataKey)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgAskTime})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventTime, userID):
		timeStr := strings.TrimSpace(userInput)
		if _, err := parseTimeLocal(timeStr); err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.ErrTimeFormat})
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgAskTime})
			return
		}
		h.addEventTimeState(userID, timeStr, data, dataKey)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgConfirm})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventDone, userID):
		if err := h.addEventDoneState(ctx, userID, data, dataKey); err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.ErrGeneric})
			return
		}
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: texts.MsgCreated})

	default:
		return
	}
}

func parseDateLocal(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", strings.TrimSpace(s), time.Local)
}

func parseTimeLocal(s string) (time.Time, error) {
	return time.ParseInLocation("15:04", strings.TrimSpace(s), time.Local)
}

func (h *Handle) getUserData(userID int64, scenarioKey string) (string, map[string]string) {
	dataKey := eventsstate.GetUserStateDataKey(scenarioKey, strconv.FormatInt(userID, 10))

	data, err := cache.GetTyped[map[string]string](h.cache, dataKey)
	if err != nil || data == nil {
		data = map[string]string{}
	}
	return dataKey, data
}

func (h *Handle) getUserState(userID int64) (string, error) {
	key := eventsstate.GetUserStateKey(userID)
	state, err := h.cache.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("get user state: %w", err)
	}
	return state, nil
}

func (h *Handle) addEventTitleState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[cacheKeyTitleValue] = userInput
	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventDate, userID), stateTTL)
	_ = cache.SetTyped(h.cache, dataKey, data, stateTTL)
}

func (h *Handle) addEventDateState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[cacheKeyDateValue] = userInput
	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventTime, userID), stateTTL)
	_ = cache.SetTyped(h.cache, dataKey, data, stateTTL)
}

func (h *Handle) addEventTimeState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[cacheKeyTimeValue] = userInput
	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventDone, userID), stateTTL)
	_ = cache.SetTyped(h.cache, dataKey, data, stateTTL)
}

func (h *Handle) addEventDoneState(ctx context.Context, userID int64, data map[string]string, dataKey string) error {
	_ = h.cache.Delete(eventsstate.GetUserStateKey(userID))
	_ = h.cache.Delete(dataKey)

	dateParsed, _ := time.Parse("2006-01-02", data[cacheKeyDateValue])
	timeParsed, _ := time.Parse("15:04", data[cacheKeyTimeValue])

	return h.eventService.AddEvent(ctx, int(userID), data[cacheKeyTitleValue], dateParsed, timeParsed)
}
