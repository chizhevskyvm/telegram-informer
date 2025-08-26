package utils

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Send(ctx context.Context, b *bot.Bot, chatID int64, text string) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	return err
}

func SendHTML(ctx context.Context, b *bot.Bot, chatID int64, text string) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
	return err
}

func AnswerOK(ctx context.Context, b *bot.Bot, update *models.Update) error {
	if update == nil || update.CallbackQuery == nil {
		return nil
	}

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	return err
}
