package events

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"telegram-informer/internal/db"
)

func HandleDeleteAllEventsToday(storage db.StorageHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := int(update.CallbackQuery.From.ID)

		err := storage.DeleteEventFromToday(ctx, userID)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "❌ Не удалось отменить события. Попробуйте позже.",
			})
			return
		}

		// Очистим кнопки на старом сообщении
		_, _ = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{},
			},
		})

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "✅ Все события на сегодня отменены.",
		})
	}
}
