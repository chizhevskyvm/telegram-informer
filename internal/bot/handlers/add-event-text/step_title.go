package addeventtext

import (
	"context"
	"strings"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/ui/texts"

	"github.com/go-telegram/bot"
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
	err = h.stateStore.SetEventAddDateState(userID)

	err = utils.SendHTML(ctx, b, chatID, texts.MsgAskDate)

	return err
}
