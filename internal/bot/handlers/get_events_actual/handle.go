package get_events_actual

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	botcommon "telegram-informer/common/bot"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/handlers"
	"telegram-informer/internal/bot/ui/texts"
	"telegram-informer/internal/domain"
)

type EventService interface {
	GetEventsActual(ctx context.Context, userId int) ([]domain.Event, error)
}

type Handle struct {
	eventService EventService
}

func NewHandle(eventService EventService) *Handle {
	return &Handle{eventService: eventService}
}

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if botcommon.BodyIsNil(update) {
		return
	}

	err := botcommon.AnswerOK(ctx, b, update)

	userID := botcommon.GetUserID(update)
	chatID := botcommon.GetChatID(update)

	events, err := h.eventService.GetEventsActual(ctx, int(userID))
	if err != nil {
		err = botcommon.SendHTML(ctx, b, chatID, texts.MsgNoEventsActual)
		return
	}

	if len(events) == 0 {
		err = botcommon.SendHTML(ctx, b, chatID, texts.MsgNoEventsActual)
		return
	}

	replyMarkup := buildEventsKeyboard(events)
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        texts.MsgEventsActualList,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: replyMarkup,
	})

	if err != nil {
		err = botcommon.Send(ctx, b, chatID, texts.ErrGeneric)
		fmt.Printf("error: %v\n", err)
	}
}

func buildEventsKeyboard(events []domain.Event) *models.InlineKeyboardMarkup {
	buttons := make([][]models.InlineKeyboardButton, 0, len(events))
	for _, e := range events {
		label := fmt.Sprintf(texts.BtnEventFormat, e.Title, utils.FormatTime(e.TimeToNotify))
		button := models.InlineKeyboardButton{
			Text:         label,
			CallbackData: fmt.Sprintf("%s%d", handlers.CBGetEventById, e.ID),
		}
		buttons = append(buttons, []models.InlineKeyboardButton{button})
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}
