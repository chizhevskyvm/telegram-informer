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

var menuKeyboard = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: texts.BtnAddEvent, CallbackData: handlers.CBAddEvent}},
		{
			{Text: texts.BtnTodayEvents, CallbackData: handlers.CBTodayEvents},
			{Text: texts.BtnAllEvents, CallbackData: handlers.CBAllEvents},
		},
		{{Text: texts.BtnCancelAllTodayEvents, CallbackData: handlers.CBCancelAllTodayEvents}},
	},
}

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        texts.BtnMenuTitle,
		ReplyMarkup: menuKeyboard,
	})

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
