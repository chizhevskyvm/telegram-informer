package events

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
	"telegram-informer/internal/db"
)

const ()

func GetId(update *models.Update) (int, error) {
	rawID := update.CallbackQuery.Data // "get-by-id:13"
	parts := strings.Split(rawID, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf(rawID)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf(rawID)
	}

	return id, nil
}

func HandleGetEventByIdCallback(storage db.StorageHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		id, _ := GetId(update)

		event, err := storage.GetEvent(ctx, int(update.CallbackQuery.From.ID), id)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "❌ Не удалось найти событие. Пожалуйста, попробуйте позже.",
			})
			return
		}

		messageText := fmt.Sprintf(
			"📅 <b>%s</b>\n🕒 Когда напомнить: <b>%s в %s</b>\n",
			event.Title,
			event.Notification.Format("02.01.2006"),
			event.TimeToNotify.Format("15:04"),
		)

		deleteButton := [][]models.InlineKeyboardButton{
			{
				{
					Text:         "🗑 Удалить",
					CallbackData: fmt.Sprintf("%s%d", CBDeleteById, event.ID),
				},
			},
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        messageText,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: deleteButton},
		})
	}
}

func HandleDeleteEventByIdCallback(storage db.StorageHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		id, _ := GetId(update)
		err := storage.DeleteEvent(ctx, int(update.CallbackQuery.From.ID), id)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "❌ Не удалось найти событие. Пожалуйста, попробуйте позже.",
			})
			return
		}
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "✅ Событие успешно удалено.",
		})
	}
}
