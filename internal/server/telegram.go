package server

import (
	"github.com/go-telegram/bot"
	redis "telegram-informer/internal/cache"
	"telegram-informer/internal/db"
	"telegram-informer/internal/handlers/tg/events"
)

const text = ""

func RegisterHandlers(b *bot.Bot, storage db.StorageHandler, cache redis.Cache) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, events.HandleStart)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBAddEvent, bot.MatchTypeExact, events.HandleAddCallback(cache))
	b.RegisterHandler(bot.HandlerTypeMessageText, text, bot.MatchTypePrefix, events.HandleAddEventText(storage, cache))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBTodayEvents, bot.MatchTypeExact, events.HandleGetEventToday(storage))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, events.CBCancelAllTodayEvents, bot.MatchTypeExact, events.HandleDeleteAllEventsToday(storage))
}
