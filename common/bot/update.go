package bot

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func AnswerOK(ctx context.Context, b *bot.Bot, update *models.Update) error {
	if update == nil || update.CallbackQuery == nil {
		return nil
	}

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	return err
}

func BodyIsNil(update *models.Update) bool {
	return update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil
}

func GetUserID(update *models.Update) int64 {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	if update.Message != nil {
		return update.Message.From.ID
	}
	return 0
}

func GetChatID(update *models.Update) int64 {
	if update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil {
		return update.CallbackQuery.Message.Message.Chat.ID
	}
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	return 0
}
