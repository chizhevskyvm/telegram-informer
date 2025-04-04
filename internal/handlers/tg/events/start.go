package events

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	CBAddEvent             = "add-event"
	CBTodayEvents          = "today-events"
	CBAllEvents            = "all-events"
	CBEditEvent            = "edit-event"
	CBDeleteEvent          = "delete-event"
	CBCancelAllTodayEvents = "cancel-today"
)

func HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Выбери опцию:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "➕ Добавить событие", CallbackData: CBAddEvent}},
				{{Text: "📅 Мои события на сегодня", CallbackData: CBTodayEvents}, {Text: "🗂 Мои события", CallbackData: CBAllEvents}},
				{{Text: "✏️ Редактировать событие", CallbackData: CBEditEvent}, {Text: "🗑 Удалить событие", CallbackData: CBDeleteEvent}},
				{{Text: "❌ Отменить все планы на сегодня", CallbackData: CBCancelAllTodayEvents}},
			},
		},
	})
	if err != nil {
		fmt.Printf("send message error: %v\n", err)
	}
}
