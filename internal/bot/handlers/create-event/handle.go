package createevent

import (
	"context"
	"telegram-informer/common/utils"
	stateh "telegram-informer/internal/bot/state"
	"telegram-informer/internal/bot/ui/texts"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handle struct {
	state *stateh.Store
}

func NewHandle(state *stateh.Store) *Handle {
	return &Handle{state: state}
}

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	err := h.state.ClearEventData(userID)

	err = h.state.SetState(userID, stateh.CreateEventState(userID))

	err = utils.SendHTML(ctx, b, chatID, texts.MsgAskTitle)

	if err != nil {
		err = utils.Send(ctx, b, chatID, texts.ErrGeneric)
		print(err) //logger
	}
}
