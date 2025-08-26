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

func (h *Handle) handleDate(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) error {
	var err error

	userID := botcommon.GetUserID(update)
	chatID := botcommon.GetChatID(update)

	date, parseErr := utils.ParseDateLocal(strings.TrimSpace(update.Message.Text))
	if parseErr != nil {
		err = botcommon.SendHTML(ctx, b, chatID, texts.ErrDateFormat)
		return err
	}

	eventData, _ := h.stateStore.GetAddEventData(userID)
	eventData.SetDate(date)

	err = h.stateStore.SetAddEventData(userID, eventData)
	err = h.stateStore.SetEventAddTimeState(userID)

	err = botcommon.SendHTML(ctx, b, chatID, texts.MsgAskTime)

	return err
}
