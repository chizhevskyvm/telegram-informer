package geteventbyid

import (
	"context"
	"fmt"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/handlers/update-helper"

	"telegram-informer/internal/bot/handlers"
	"telegram-informer/internal/bot/ui/texts"
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

func NewHandle(eventService EventService) *Handler { return &Handler{eventService: eventService} }

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	err := utils.AnswerOK(ctx, b, update)

	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	id, err := updatehelper.GetId(update)
	if err != nil {
		err = utils.SendHTML(ctx, b, chatID, texts.MsgEventNotFound)

		fmt.Printf("error: %v\n", err)
		return
	}

	event, err := h.eventService.GetEvent(ctx, userID, id)
	if err != nil {
		err = utils.SendHTML(ctx, b, chatID, texts.MsgEventNotFound)

		fmt.Printf("error: %v\n", err)
		return
	}

	messageText := fmt.Sprintf(
		texts.MsgEventDetails,
		event.Title,
		event.Notification.Format("02.01.2006"),
		event.TimeToNotify.Format("15:04"),
	)

	replyMarkup := buildDeleteKeyboard(event.ID)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        messageText,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: replyMarkup,
	})

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func buildDeleteKeyboard(eventID int) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         texts.BtnDeleteEvent,
					CallbackData: fmt.Sprintf("%s%d", handlers.CBDeleteById, eventID),
				},
			},
		},
	}
}
