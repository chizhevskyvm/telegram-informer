package addeventtext

import (
	"context"
	"strings"
	botcommon "telegram-informer/common/bot"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/ui/texts"

	"github.com/go-telegram/bot/models"

	"github.com/go-telegram/bot"
)

func (h *Handle) handleTime(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) error {
	var err error

	userID := botcommon.GetUserID(update)
	chatID := botcommon.GetChatID(update)

	date, parseErr := utils.ParseTime(strings.TrimSpace(update.Message.Text))
	if parseErr != nil {
		err = botcommon.SendHTML(ctx, b, chatID, texts.ErrTimeFormat)
		return err
	}

	eventData, _ := h.stateStore.GetAddEventData(userID)
	eventData.SetTime(date)

	err = h.stateStore.SetAddEventData(userID, eventData)
	err = h.stateStore.SetDoneState(userID)

	err = botcommon.SendHTML(ctx, b, chatID, texts.MsgConfirm)

	return err
}
