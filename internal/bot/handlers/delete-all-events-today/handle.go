package deletealleventstoday

import (
	"context"

	"telegram-informer/common/utils"
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
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	err := utils.AnswerOK(ctx, b, update)

	userID := int(update.CallbackQuery.From.ID)
	msg := update.CallbackQuery.Message.Message
	chatID := msg.Chat.ID
	messageID := msg.ID

	if err = h.eventService.DeleteEventFromToday(ctx, userID); err != nil {
		_ = utils.SendHTML(ctx, b, chatID, texts.MsgDeleteAllError)
		return
	}

	empty := &models.InlineKeyboardMarkup{InlineKeyboard: make([][]models.InlineKeyboardButton, 0)}
	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: empty},
	)

	err = utils.SendHTML(ctx, b, chatID, texts.MsgDeleteAllSuccess)

	if err != nil {
		err = utils.Send(ctx, b, chatID, texts.ErrGeneric)
		print(err) //logger
	}
}
