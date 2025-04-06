package events

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"telegram-informer/internal/db"
)

func HandleGetEventToday(storage db.StorageHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		events, err := storage.GetEventsTodayFromUser(ctx, int(update.CallbackQuery.From.ID))
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "❌ Пожалуйста, попробуйте позже. "})
		}

		if len(events) == 0 {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "\"📅 Всё чисто! Можно валяться весь день 😎\"",
			})

			return
		}

		buttons := make([][]models.InlineKeyboardButton, 0, len(events))
		for _, e := range events {
			label := fmt.Sprintf("📌 %s — 🕒 %s", e.Title, e.TimeToNotify.Format("15:04"))
			button := models.InlineKeyboardButton{
				Text:         label,
				CallbackData: fmt.Sprintf("%s%d", CBGetById, e.ID),
			}
			buttons = append(buttons, []models.InlineKeyboardButton{button})
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "📅 События на сегодня: ",
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: buttons,
			},
		})
	}
}
