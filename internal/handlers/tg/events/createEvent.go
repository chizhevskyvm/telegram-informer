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
		state, err := cache.Get(tg.GetUserStateKey(userId))
		if err != nil && !errors.Is(err, redis.Nil) {
			fmt.Println("Ошибка при получении состояния:", err)
			return
		}

		if isAddEventState(state, userId) {
			state = tg.GetUserStageState(StateAddEventTitle, userId)
			_ = cache.Set(tg.GetUserStateKey(userId), state, time.Minute*10)
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "✏️ Введите заголовок события.\nНапример: \"День рождения бабушки 🎉\""})
	}
}

func HandleAddEventText(storage db.StorageHandler, cache iRedis.Cache) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		userInput := update.Message.Text
		userId := update.Message.From.ID

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
			handleAddEventTitleState(ctx, cache, b, update, userInput, data, dataKey)

		case tg.GetUserStageState(StateAddEventDate, userId):
			handleAddEventDateState(ctx, cache, b, update, userInput, data, dataKey)

		case tg.GetUserStageState(StateAddEventTime, userId):
			handleAddEventTimeState(ctx, cache, b, update, userInput, data, dataKey)

		case tg.GetUserStageState(StateAddEventDone, userId):
			handleAddEventDoneState(ctx, storage, cache, b, update, data, dataKey)
		}
	}
}

func handleAddEventTitleState(ctx context.Context, cache iRedis.Cache, b *bot.Bot, update *models.Update, userInput string, data map[string]string, dataKey string) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	data[titleValue] = userInput
	_ = cache.Set(tg.GetUserStateKey(userId), tg.GetUserStageState(StateAddEventDate, userId), time.Minute*10)
	_ = iRedis.SetTyped(cache, dataKey, data, time.Minute*10)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "📅 Введите дату события в формате ГГГГ-ММ-ДД.\nНапример: \"2025-12-31\""})
}

func handleAddEventDateState(ctx context.Context, cache iRedis.Cache, b *bot.Bot, update *models.Update, userInput string, data map[string]string, dataKey string) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	data[dateValue] = userInput
	_ = cache.Set(tg.GetUserStateKey(userId), tg.GetUserStageState(StateAddEventTime, userId), time.Minute*10)
	_ = iRedis.SetTyped(cache, dataKey, data, time.Minute*10)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "⏰ Введите время события в 24-часовом формате ЧЧ:ММ.\nНапример: \"14:30\""})
}

func handleAddEventTimeState(ctx context.Context, cache iRedis.Cache, b *bot.Bot, update *models.Update, userInput string, data map[string]string, dataKey string) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	data[timeValue] = userInput
	_ = cache.Set(tg.GetUserStateKey(userId), tg.GetUserStageState(StateAddEventDone, userId), time.Minute*10)
	_ = iRedis.SetTyped(cache, dataKey, data, time.Minute*10)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "✅ Всё готово! Подтвердите создание события, написав \"да\" или \"нет\"."})
}

func handleAddEventDoneState(ctx context.Context, storage db.StorageHandler, cache iRedis.Cache, b *bot.Bot, update *models.Update, data map[string]string, dataKey string) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	_ = cache.Delete(tg.GetUserStateKey(userId))
	_ = cache.Delete(dataKey)

	dateParsed, _ := time.Parse("2006-01-02", data[dateValue])
	timeParsed, _ := time.Parse("15:04", data[timeValue])

	_ = storage.AddEvent(ctx, int(userId), data[titleValue], dateParsed, timeParsed)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "🎉 Событие создано!",
	})
}

func isAddEventState(state string, userId int64) bool {
	return state != tg.GetUserStageState(StateAddEventTitle, userId) &&
		state != tg.GetUserStageState(StateAddEventDate, userId) &&
		state != tg.GetUserStageState(StateAddEventTime, userId) &&
		state != tg.GetUserStageState(StateAddEventDone, userId)
}
