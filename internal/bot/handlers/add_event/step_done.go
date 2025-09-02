package addeventtext

import (
	"context"
	"fmt"
	"strings"
	botcommon "telegram-informer/common/bot"
	"telegram-informer/common/utils"
	"telegram-informer/internal/bot/ui/texts"
	"telegram-informer/internal/domain"

	"github.com/go-telegram/bot/models"

	"github.com/go-telegram/bot"
)

func (h *Handle) handleDone(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userID := botcommon.GetUserID(update)
	chatID := botcommon.GetChatID(update)

	answer := strings.ToLower(strings.TrimSpace(update.Message.Text))

	var err error

	switch {
	case isYes(answer):
		ed, _ := h.stateStore.GetAddEventData(userID)
		if err = h.addEventDone(ctx, userID, ed); err != nil {
			err = botcommon.SendHTML(ctx, b, chatID, texts.ErrGeneric)
			return nil
		}
		err = h.stateStore.ClearEventData(userID)
		err = h.stateStore.ClearState(userID)

		if err == nil {
			err = SendEventCreatedDetails(ctx, b, chatID, ed)
		}
	case isNo(answer):
		err = h.stateStore.SetCreateEventState(userID)
		err = botcommon.SendHTML(ctx, b, chatID, texts.MsgAskTitle)
	default:
		err = botcommon.SendHTML(ctx, b, chatID, texts.ErrYesOrNo)
		err = botcommon.SendHTML(ctx, b, chatID, texts.MsgConfirm)
	}

	return err
}

func (h *Handle) addEventDone(ctx context.Context, userID int64, ed *domain.EventData) error {
	d, _ := ed.GetDate()
	t, _ := ed.GetTime()
	return h.eventService.AddEvent(ctx, int(userID), ed.GetTitle(), d, t)
}

func SendEventCreatedDetails(ctx context.Context, b *bot.Bot, chatID int64, ed *domain.EventData) error {
	date, _ := ed.GetDate()
	time, _ := ed.GetTime()
	msg := fmt.Sprintf(
		texts.MsgCreated,
		ed.GetTitle(),
		utils.FormatDate(date),
		utils.FormatTime(time),
	)
	return botcommon.SendHTML(ctx, b, chatID, msg)
}

func isYes(s string) bool {
	switch s {
	case "да", "д", "yes", "y", "+", "ок", "okay", "ага":
		return true
	}
	return false
}

func isNo(s string) bool {
	switch s {
	case "нет", "не", "no", "n", "-", "nope":
		return true
	}
	return false
}
