package start

import (
	"context"
	"fmt"
	"telegram-informer/internal/bot/handlers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
}

func NewHandle() *Handler {
	return &Handler{}
}

func (h Handler) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Меню: ",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "➕ Добавить событие", CallbackData: handlers.CBAddEvent}},
				{{Text: "📅 Мои события на сегодня", CallbackData: handlers.CBTodayEvents}, {Text: "🗂 Мои события", CallbackData: handlers.CBAllEvents}},
				{{Text: "❌ Отменить все планы на сегодня", CallbackData: handlers.CBCancelAllTodayEvents}},
			},
		},
	})
	if err != nil {
		fmt.Printf("send message error: %v\n", err)
	}
}
