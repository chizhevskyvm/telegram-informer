package addeventtext

import (
	"context"
	"strings"
	botcommon "telegram-informer/common/bot"
	"telegram-informer/internal/bot/ui/texts"

	"github.com/go-telegram/bot/models"

	"github.com/go-telegram/bot"
)

func (h *Handle) handleTitle(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) error {
	userID := botcommon.GetUserID(update)
	chatID := botcommon.GetChatID(update)

	title := strings.TrimSpace(update.Message.Text)
	if title == "" {
		_ = botcommon.SendHTML(ctx, b, chatID, texts.ErrTitleEmpty)
		return nil
	}

	eventData, _ := h.stateStore.GetAddEventData(userID)
	eventData.SetTitle(title)

	err := h.stateStore.SetAddEventData(userID, eventData)
	err = h.stateStore.SetEventAddDateState(userID)

	err = botcommon.SendHTML(ctx, b, chatID, texts.MsgAskDate)

	return err
}
