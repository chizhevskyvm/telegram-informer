package addeventtext

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	eventsstate "telegram-informer/internal/bot/event-state"
	"telegram-informer/internal/infra/cache"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
)

const (
	titleValue = "title_value"
	dateValue  = "date_value"
	timeValue  = "time_value"
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

func (h Handle) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userInput := update.Message.Text
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	state, err := h.cache.Get(eventsstate.GetUserStateKey(userId))
	if err != nil && !errors.Is(err, redis.Nil) {
		fmt.Println("Ошибка при получении состояния: ", err)
		return
	}

	if eventsstate.IsAddEventState(state, userId) {
		return
	}

	dataKey := eventsstate.GetUserStateDataKey("addEvent", strconv.FormatInt(userId, 10))
	data, err := cache.GetTyped[map[string]string](h.cache, dataKey)
	if err != nil {
		data = map[string]string{}
	}

	switch state {
	case eventsstate.GetUserStageState(eventsstate.StateAddEventTitle, userId):
		h.addEventTitleState(userId, userInput, data, dataKey)

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "📅 Введите дату события в формате ГГГГ-ММ-ДД.\nНапример: \"2025-12-31\""})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventDate, userId):
		h.addEventDateState(userId, userInput, data, dataKey)

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "⏰ Введите время события в 24-часовом формате ЧЧ:ММ.\nНапример: \"14:30\""})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventTime, userId):

		h.addEventTimeState(update, userInput, data, dataKey)

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "✅ Всё готово! Подтвердите создание события, написав \"да\" или \"нет\"."})

	case eventsstate.GetUserStageState(eventsstate.StateAddEventDone, userId):
		h.addEventDoneState(ctx, userId, data, dataKey)

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "🎉 Событие создано!",
		})
	}
}

func (h Handle) addEventTitleState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[titleValue] = userInput

	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventDate, userID), time.Minute*10)
	_ = cache.SetTyped(h.cache, dataKey, data, time.Minute*10)
}

func (h Handle) addEventDateState(userID int64, userInput string, data map[string]string, dataKey string) {
	data[dateValue] = userInput

	_ = h.cache.Set(eventsstate.GetUserStateKey(userID), eventsstate.GetUserStageState(eventsstate.StateAddEventTime, userID), time.Minute*10)
	_ = cache.SetTyped(h.cache, dataKey, data, time.Minute*10)
}

func (h Handle) addEventTimeState(update *models.Update, userInput string, data map[string]string, dataKey string) {
	userId := update.Message.From.ID
	data[timeValue] = userInput
	_ = h.cache.Set(eventsstate.GetUserStateKey(userId), eventsstate.GetUserStageState(eventsstate.StateAddEventDone, userId), time.Minute*10)
	_ = cache.SetTyped(h.cache, dataKey, data, time.Minute*10)

}

func (h Handle) addEventDoneState(ctx context.Context, userID int64, data map[string]string, dataKey string) {
	_ = h.cache.Delete(eventsstate.GetUserStateKey(userID))
	_ = h.cache.Delete(dataKey)

	dateParsed, _ := time.Parse("2006-01-02", data[dateValue])
	timeParsed, _ := time.Parse("15:04", data[timeValue])

	_ = h.eventService.AddEvent(ctx, int(userID), data[titleValue], dateParsed, timeParsed)
}
