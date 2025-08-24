package deleteeventbyid

import (
	"context"
	updatehelper "telegram-informer/internal/bot/update-helper"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type EventService interface {
	DeleteEvent(ctx context.Context, userId int, eventId int) error
}

type Handle struct {
	eventService EventService
}

func NewHandle(eventService EventService) Handle {
	return Handle{eventService: eventService}
}

func (h Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	id, _ := updatehelper.GetId(update)
	err := h.eventService.DeleteEvent(ctx, int(update.CallbackQuery.From.ID), id)
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
