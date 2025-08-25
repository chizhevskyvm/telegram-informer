package geteventbyid

import (
	"context"
	"fmt"
	"telegram-informer/internal/bot/handlers"
	"telegram-informer/internal/bot/ui/texts"
	updatehelper "telegram-informer/internal/bot/update-helper"
	"telegram-informer/internal/domain"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type EventService interface {
	DeleteEvent(ctx context.Context, userId int, eventId int) error
	GetEvent(ctx context.Context, userId int, id int) (domain.Event, error)
}

type Handler struct {
	eventService EventService
}

func NewHandle(eventService EventService) *Handler {
	return &Handler{eventService: eventService}
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	id, _ := updatehelper.GetId(update)
	event, err := h.eventService.GetEvent(ctx, userID, id)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   texts.MsgEventNotFound,
		})
		return
	}

	messageText := fmt.Sprintf(
		texts.MsgEventDetails,
		event.Title,
		event.Notification.Format("02.01.2006"),
		event.TimeToNotify.Format("15:04"),
	)

	deleteButton := [][]models.InlineKeyboardButton{
		{
			{
				Text:         texts.BtnDeleteEvent,
				CallbackData: fmt.Sprintf("%s%d", handlers.CBDeleteById, event.ID),
			},
		},
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        messageText,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: deleteButton},
	})
}
