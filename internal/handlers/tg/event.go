package tg

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"telegram-informer/internal/db"
)

func HandleAddEvent(storage db.StorageHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
	}
}

func HandleGetEvent(storage db.StorageHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

	}
}
func HandleDeleteEvent(storage db.StorageHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

	}
}
