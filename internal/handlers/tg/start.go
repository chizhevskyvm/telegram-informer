package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Выбери опцию:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "Добавить новое событие", CallbackData: "add-event"}},
				{{Text: "Получить события на сегодня", CallbackData: "get-event-today"}},
				{{Text: "Удалить событие", CallbackData: "delete-event"}},
			},
		},
	})

	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}
