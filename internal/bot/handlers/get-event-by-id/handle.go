package geteventbyid

import (
	"context"
	"fmt"
	botcommon "telegram-informer/common/bot"
	"telegram-informer/common/utils"
	updatehelper "telegram-informer/internal/bot/handlers/update-helper"

	"github.com/go-telegram/bot"

	"telegram-informer/internal/bot/handlers"
	"telegram-informer/internal/bot/ui/texts"
	"telegram-informer/internal/domain"

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
	if botcommon.BodyIsNil(update) {
		return
	}

	err := botcommon.AnswerOK(ctx, b, update)

	userID := int(botcommon.GetUserID(update))
	chatID := botcommon.GetChatID(update)

	id, err := updatehelper.ParseCallbackID(update)
	if err != nil {
		err = botcommon.SendHTML(ctx, b, chatID, texts.MsgEventNotFound)

		fmt.Printf("error: %v\n", err)
		return
	}

	event, err := h.eventService.GetEvent(ctx, userID, id)
	if err != nil {
		err = botcommon.SendHTML(ctx, b, chatID, texts.MsgEventNotFound)

		fmt.Printf("error: %v\n", err)
		return
	}

	messageText := fmt.Sprintf(
		texts.MsgEventDetails,
		event.Title,
		utils.FormatDate(event.Notification),
		utils.FormatTime(event.TimeToNotify),
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
					CallbackData: fmt.Sprintf("%s%d", handlers.CBDeleteEventById, eventID),
				},
			},
		},
	}
}
