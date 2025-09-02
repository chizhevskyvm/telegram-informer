package deleteeventbyid

import (
	"context"
	botcommon "telegram-informer/common/bot"
	updatehelper "telegram-informer/internal/bot/handlers/update_helper"
	"telegram-informer/internal/bot/ui/texts"

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

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if botcommon.BodyIsNil(update) {
		return
	}

	err := botcommon.AnswerOK(ctx, b, update)

	userID := int(update.CallbackQuery.From.ID)
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	id, err := updatehelper.ParseCallbackID(update)
	if err != nil {
		_ = botcommon.SendHTML(ctx, b, chatID, texts.MsgDeleteError)
		return
	}

	if err = h.eventService.DeleteEvent(ctx, userID, id); err != nil {
		err = botcommon.SendHTML(ctx, b, chatID, texts.MsgDeleteError)
		return
	}

	err = botcommon.SendHTML(ctx, b, chatID, texts.MsgDeleteSuccess)

	if err != nil {
		err = botcommon.Send(ctx, b, chatID, texts.ErrGeneric)
		print(err) //logger
	}
}
