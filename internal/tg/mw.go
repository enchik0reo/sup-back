package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AdminOnly(adminSet map[string]struct{}, next ViewFunc) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		if _, ok := adminSet[update.Message.From.UserName]; ok {
			return next(ctx, bot, update)
		}

		if _, err := bot.api.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"У вас нет прав.",
		)); err != nil {
			return err
		}

		return nil
	}
}
