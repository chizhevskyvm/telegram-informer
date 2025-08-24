package geteventtoday

import (
	"context"
	"fmt"
	"telegram-informer/internal/bot/handlers"
	"telegram-informer/internal/domain"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type EventService interface {
	GetEventsTodayFromUser(ctx context.Context, userId int) ([]domain.Event, error)
}

type Handle struct {
	eventService EventService
}

func NewHandle(eventService EventService) *Handle {
	return &Handle{eventService: eventService}
}

func (h Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	events, err := h.eventService.GetEventsTodayFromUser(ctx, int(update.CallbackQuery.From.ID))
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
			CallbackData: fmt.Sprintf("%s%d", handlers.CBGetById, e.ID),
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
