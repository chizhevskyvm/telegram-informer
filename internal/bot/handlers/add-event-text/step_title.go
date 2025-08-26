package addeventtext

import (
	"context"
	"github.com/go-telegram/bot"
	"strings"
	"telegram-informer/common/utils"
	stateh "telegram-informer/internal/bot/state"

	"telegram-informer/internal/bot/ui/texts"
)

func (h *Handle) handleTitle(ctx context.Context, b *bot.Bot, chatID int64, userID int64, userInput string) error {
	title := strings.TrimSpace(userInput)
	if title == "" {
		_ = utils.SendHTML(ctx, b, chatID, texts.ErrTitleEmpty)
		return nil
	}

	eventData, _ := h.stateStore.GetAddEventData(userID)
	eventData.SetTitle(title)

	err := h.stateStore.SetAddEventData(userID, eventData)
	err = h.stateStore.SetState(userID, stateh.AddEventDateState(userID))

	err = utils.SendHTML(ctx, b, chatID, texts.MsgAskDate)

	return err
}
