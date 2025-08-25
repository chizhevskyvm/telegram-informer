package mainmenu

import (
	"context"
	"fmt"
	"telegram-informer/internal/bot/ui/texts"

	"telegram-informer/internal/bot/handlers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handle struct{}

func NewHandle() *Handle {
	return &Handle{}
}

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   texts.BtnMenuTitle,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: texts.BtnAddEvent, CallbackData: handlers.CBAddEvent}},
				{{Text: texts.BtnTodayEvents, CallbackData: handlers.CBTodayEvents}, {Text: texts.BtnAllEvents, CallbackData: handlers.CBAllEvents}},
				{{Text: texts.BtnCancelAllTodayEvents, CallbackData: handlers.CBCancelAllTodayEvents}},
			},
		},
	})
	if err != nil {
		fmt.Printf("send message error: %v\n", err)
	}
}
