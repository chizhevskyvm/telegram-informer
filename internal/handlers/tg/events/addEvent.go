package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
	"strconv"
	iRedis "telegram-informer/internal/cache"
	"telegram-informer/internal/db"
	"telegram-informer/internal/handlers/tg"
	"time"
)

const (
	StateAddEventTitle = "add_event:%d:title"
	StateAddEventDate  = "add_event:%d:date"
	StateAddEventTime  = "add_event:%d:time"
	StateAddEventDone  = "add_event:%d:done"

	titleValue = "title_value"
	dateValue  = "date_value"
	timeValue  = "time_value"
)

func HandleAddCallback(cache iRedis.Cache) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.CallbackQuery == nil {
			return
		}

		userId := update.CallbackQuery.From.ID
		chatId := update.CallbackQuery.Message.Message.Chat.ID

		state, err := cache.Get(tg.GetUserStateKey(userId))
		if err != nil && !errors.Is(err, redis.Nil) {
			fmt.Println("Ошибка при получении состояния:", err)
			return
		}

		if isAddEventState(state, userId) {
			state = tg.GetUserStageState(StateAddEventTitle, userId)
			_ = cache.Set(tg.GetUserStateKey(userId), state, time.Minute*10)
		}

		handleTitleState(ctx, b, chatId)
	}
}

func HandleAddEventText(storage db.StorageHandler, cache iRedis.Cache) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		userId := update.Message.From.ID
		chatId := update.Message.Chat.ID
		userInput := update.Message.Text

		state, err := cache.Get(tg.GetUserStateKey(userId))
		if err != nil && !errors.Is(err, redis.Nil) {
			fmt.Println("Ошибка при получении состояния: ", err)
			return
		}

		if isAddEventState(state, userId) {
			return
		}

		dataKey := tg.GetUserStateDataKey("addEvent", strconv.FormatInt(userId, 10))
		data, err := iRedis.GetTyped[map[string]string](cache, dataKey)
		if err != nil {
			data = map[string]string{}
		}

		switch state {
		case tg.GetUserStageState(StateAddEventTitle, userId):
			handleAddEventTitleState(ctx, cache, b, userId, chatId, userInput, data, dataKey)

		case tg.GetUserStageState(StateAddEventDate, userId):
			handleAddEventDateState(ctx, cache, b, userId, chatId, userInput, data, dataKey)

		case tg.GetUserStageState(StateAddEventTime, userId):
			handleAddEventTimeState(ctx, cache, b, userId, chatId, userInput, data, dataKey)

		case tg.GetUserStageState(StateAddEventDone, userId):
			handleAddEventDoneState(ctx, storage, cache, b, userId, chatId, data, dataKey)
		}
	}
}

func handleAddEventTitleState(ctx context.Context, cache iRedis.Cache, b *bot.Bot, userId int64, chatId int64, userInput string, data map[string]string, dataKey string) {
	data[titleValue] = userInput
	_ = cache.Set(tg.GetUserStateKey(userId), tg.GetUserStageState(StateAddEventDate, userId), time.Minute*10)
	_ = iRedis.SetTyped(cache, dataKey, data, time.Minute*10)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "📅 Введите дату события в формате ГГГГ-ММ-ДД.\nНапример: \"2025-12-31\""})
}

func handleAddEventDateState(ctx context.Context, cache iRedis.Cache, b *bot.Bot, userId int64, chatId int64, userInput string, data map[string]string, dataKey string) {
	data[dateValue] = userInput
	_ = cache.Set(tg.GetUserStateKey(userId), tg.GetUserStageState(StateAddEventTime, userId), time.Minute*10)
	_ = iRedis.SetTyped(cache, dataKey, data, time.Minute*10)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "⏰ Введите время события в 24-часовом формате ЧЧ:ММ.\nНапример: \"14:30\""})
}

func handleAddEventTimeState(ctx context.Context, cache iRedis.Cache, b *bot.Bot, userId int64, chatId int64, userInput string, data map[string]string, dataKey string) {
	data[timeValue] = userInput
	_ = cache.Set(tg.GetUserStateKey(userId), tg.GetUserStageState(StateAddEventDone, userId), time.Minute*10)
	_ = iRedis.SetTyped(cache, dataKey, data, time.Minute*10)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "✅ Всё готово! Подтвердите создание события, написав \"да\" или \"нет\"."})
}

func handleAddEventDoneState(ctx context.Context, storage db.StorageHandler, cache iRedis.Cache, b *bot.Bot, userId int64, chatId int64, data map[string]string, dataKey string) {
	_ = cache.Delete(tg.GetUserStateKey(userId))
	_ = cache.Delete(dataKey)

	title := data[titleValue]
	date := data[dateValue]
	timeStr := data[timeValue]

	datetimeStr := fmt.Sprintf("%sT%s:00", date, timeStr)
	eventTime, err := time.Parse("2006-01-02T15:04:05", datetimeStr)
	if err != nil {
		fmt.Printf("Ошибка парсинга даты и времени: %v\n", err)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "❌ Ошибка при создании события. Убедитесь, что дата и время введены верно.",
		})
		return
	}

	_ = storage.AddEvent(ctx, int(userId), title, eventTime, eventTime)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "🎉 Событие создано!",
	})
}

func handleTitleState(ctx context.Context, b *bot.Bot, chatID int64) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "✏️ Введите заголовок события.\nНапример: \"День рождения бабушки 🎉\""})
}

func isAddEventState(state string, userId int64) bool {
	return state != tg.GetUserStageState(StateAddEventTitle, userId) &&
		state != tg.GetUserStageState(StateAddEventDate, userId) &&
		state != tg.GetUserStageState(StateAddEventTime, userId) &&
		state != tg.GetUserStageState(StateAddEventDone, userId)
}
