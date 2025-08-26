package addeventtext

import (
	"context"
	"github.com/go-telegram/bot"
	"strings"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/ui/texts"
)

func (h *Handle) handleDate(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	userID int64,
	userInput string,
) error {
	var err error

	date, parseErr := utils.ParseDateLocal(strings.TrimSpace(userInput))
	if parseErr != nil {
		err = utils.SendHTML(ctx, b, chatID, texts.ErrDateFormat)
		return err
	}

	eventData, _ := h.stateStore.GetAddEventData(userID)
	eventData.SetDate(date)

	err = h.stateStore.SetAddEventData(userID, eventData)
	err = h.stateStore.SetEventAddTimeState(userID)

	err = utils.SendHTML(ctx, b, chatID, texts.MsgAskTime)

	return err
}
