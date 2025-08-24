package createevent

import (
	"context"
	"errors"
	"fmt"
	eventsstate "telegram-informer/internal/bot/event-state"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
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

func (h Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	userId := update.CallbackQuery.From.ID
	state, err := h.cache.Get(eventsstate.GetUserStateKey(userId))
	if err != nil && !errors.Is(err, redis.Nil) {
		fmt.Println("Ошибка при получении состояния:", err)
		return
	}

	if eventsstate.IsAddEventState(state, userId) {
		state = eventsstate.GetUserStageState(eventsstate.StateAddEventTitle, userId)
		_ = h.cache.Set(eventsstate.GetUserStateKey(userId), state, time.Minute*10)
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "✏️ Введите заголовок события.\nНапример: \"День рождения бабушки 🎉\""})
}
