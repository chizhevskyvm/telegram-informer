package createevent

import (
	"context"
	"errors"
	"fmt"
	eventsstate "telegram-informer/internal/bot/event-state"
	"telegram-informer/internal/bot/ui/texts"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
)

// UI-константы
const (
	msgEnterTitle = "✏️ Введите заголовок события.\nНапример: \"День рождения бабушки 🎉\""
	stateTTL      = 10 * time.Minute
)

type Cache interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Handle struct {
	cache Cache
}

func NewHandle(cache Cache) *Handle {
	return &Handle{cache: cache}
}

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	state, err := h.cache.Get(eventsstate.GetUserStateKey(userID))
	if err != nil && !errors.Is(err, redis.Nil) {
		fmt.Println("Ошибка при получении состояния:", err)
		return
	}

	if eventsstate.IsAddEventState(state, userID) {
		state = eventsstate.GetUserStageState(eventsstate.StateAddEventTitle, userID)
		_ = h.cache.Set(eventsstate.GetUserStateKey(userID), state, stateTTL)
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   texts.MsgAskTitle,
	})
}
