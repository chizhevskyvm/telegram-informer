package addeventtext

import (
	"context"
	"fmt"
	botcommon "telegram-informer/common/bot"
	"telegram-informer/internal/bot/ui/texts"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	stateh "telegram-informer/internal/bot/state"
	"time"
)

type Handle struct {
	eventService EventService
	stateStore   *stateh.Store
}

func NewHandle(eventService EventService, st *stateh.Store) *Handle {
	return &Handle{eventService: eventService, stateStore: st}
}

type EventService interface {
	AddEvent(ctx context.Context, userId int, title string, time time.Time, timeToNotify time.Time) error
}

func (h *Handle) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := botcommon.GetUserID(update)
	chatID := botcommon.GetChatID(update)

	state, err := h.stateStore.GetState(userID)
	if err != nil {
		fmt.Println("Ошибка при получении состояния:", err)
		return
	}
	if !(stateh.IsAddEventState(state, userID) || stateh.IsCreateEventState(state, userID)) {
		return
	}

	if state == stateh.CreateEventState(userID) {
		err = h.stateStore.SetEventAddTitleState(userID)
	}

	switch state {
	case stateh.AddEventTitleState(userID):
		err = h.handleTitle(ctx, b, update)
	case stateh.AddEventDateState(userID):
		err = h.handleDate(ctx, b, update)
	case stateh.AddEventTimeState(userID):
		err = h.handleTime(ctx, b, update)
	case stateh.AddEventDoneState(userID):
		err = h.handleDone(ctx, b, update)
	default:
		return
	}

	if err != nil {
		err = botcommon.Send(ctx, b, chatID, texts.ErrGeneric)
		print(err) //logger
	}
}
