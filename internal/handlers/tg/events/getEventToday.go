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
		events, err := storage.GetEvents(ctx, int(update.CallbackQuery.From.ID))
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "Видимо ничего нет :("})
		}

		var buttons [][]models.InlineKeyboardButton
		for _, event := range events {
			text := fmt.Sprintf("📌 %s (%s)", event.Title, event.TimeToNotify.Format("15:04"))
			buttons = append(buttons, []models.InlineKeyboardButton{
				{
					Text:         text,
					CallbackData: "noop",
				},
			})
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "События на сегодня:",
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: buttons,
			},
		})
	}
}
