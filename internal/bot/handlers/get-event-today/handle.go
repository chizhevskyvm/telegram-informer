package geteventtoday

import (
	"context"
	"fmt"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/handlers"
	"telegram-informer/internal/bot/ui/texts"
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

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	err := utils.AnswerOK(ctx, b, update)

	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	events, err := h.eventService.GetEventsTodayFromUser(ctx, userID)
	if err != nil {
		err = utils.SendHTML(ctx, b, chatID, texts.MsgNoEventsToday)
		return
	}

	if len(events) == 0 {
		err = utils.SendHTML(ctx, b, chatID, texts.MsgNoEventsToday)
		return
	}

	replyMarkup := buildEventsKeyboard(events)
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        texts.MsgEventsList,
		ParseMode:   models.ParseModeHTML, // у тебя тексты с HTML
		ReplyMarkup: replyMarkup,
	})

	if err != nil {
		err = utils.Send(ctx, b, chatID, texts.ErrGeneric)
		fmt.Printf("error: %v\n", err)
	}
}

func buildEventsKeyboard(events []domain.Event) *models.InlineKeyboardMarkup {
	buttons := make([][]models.InlineKeyboardButton, 0, len(events))
	for _, e := range events {
		label := fmt.Sprintf(texts.BtnEventFormat, e.Title, e.TimeToNotify.Format("15:04"))
		button := models.InlineKeyboardButton{
			Text:         label,
			CallbackData: fmt.Sprintf("%s%d", handlers.CBGetById, e.ID),
		}
		buttons = append(buttons, []models.InlineKeyboardButton{button})
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}
