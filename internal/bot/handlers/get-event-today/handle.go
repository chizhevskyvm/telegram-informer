package geteventtoday

import (
	"context"
	"fmt"
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
	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	events, err := h.eventService.GetEventsTodayFromUser(ctx, userID)
	if err != nil || len(events) == 0 {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   texts.MsgNoEventsToday,
		})
		return
	}

	buttons := make([][]models.InlineKeyboardButton, 0, len(events))
	for _, e := range events {
		label := fmt.Sprintf(texts.BtnEventFormat, e.Title, e.TimeToNotify.Format("15:04"))
		button := models.InlineKeyboardButton{
			Text:         label,
			CallbackData: fmt.Sprintf("%s%d", handlers.CBGetById, e.ID),
		}
		buttons = append(buttons, []models.InlineKeyboardButton{button})
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   texts.MsgEventsList,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	})
}
