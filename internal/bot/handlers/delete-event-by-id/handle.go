package deleteeventbyid

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/handlers/update-helper"
	"telegram-informer/internal/bot/ui/texts"
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

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	err := utils.AnswerOK(ctx, b, update)

	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	id, err := updatehelper.GetId(update)
	if err != nil {
		_ = utils.SendHTML(ctx, b, chatID, texts.MsgDeleteError)
		return
	}

	if err = h.eventService.DeleteEvent(ctx, userID, id); err != nil {
		err = utils.SendHTML(ctx, b, chatID, texts.MsgDeleteError)
		return
	}

	err = utils.SendHTML(ctx, b, chatID, texts.MsgDeleteSuccess)

	if err != nil {
		err = utils.Send(ctx, b, chatID, texts.ErrGeneric)
		print(err) //logger
	}
}
