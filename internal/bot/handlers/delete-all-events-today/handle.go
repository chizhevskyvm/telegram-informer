package deletealleventstoday

import (
	"context"
	"telegram-informer/internal/bot/ui/texts"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type EventService interface {
	DeleteEventFromToday(ctx context.Context, userId int) error
}

type Handle struct {
	eventService EventService
}

func NewHandle(eventService EventService) *Handle {
	return &Handle{eventService: eventService}
}

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID

	if err := h.eventService.DeleteEventFromToday(ctx, userID); err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   texts.MsgDeleteAllError,
		})
		return
	}

	_, _ = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{}},
	})

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   texts.MsgDeleteAllSuccess,
	})
}
