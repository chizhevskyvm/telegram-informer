package server

import (
	"github.com/go-telegram/bot"
	"telegram-informer/internal/db"
	"telegram-informer/internal/handlers/tg"
)

func RegisterHandlers(b *bot.Bot, storage db.StorageHandler) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tg.HandleStart)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add-event", bot.MatchTypeExact, tg.HandleAddEvent(storage))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "get-event-today", bot.MatchTypeExact, tg.HandleGetEvent(storage))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "delete-event", bot.MatchTypeExact, tg.HandleDeleteEvent(storage))
}
