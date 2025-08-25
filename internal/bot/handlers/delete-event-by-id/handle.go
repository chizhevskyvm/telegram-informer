package deleteeventbyid

import (
	"context"
	"telegram-informer/internal/bot/ui/texts"
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
	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	id, _ := updatehelper.GetId(update)
	if err := h.eventService.DeleteEvent(ctx, userID, id); err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   texts.MsgDeleteError,
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   texts.MsgDeleteSuccess,
	})
}
