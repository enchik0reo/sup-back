package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func adminOnly(adminSet map[string]int64, next ViewFunc) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		if update.Message != nil {
			if _, ok := adminSet[update.Message.From.UserName]; ok {
				adminSet[update.Message.From.UserName] = update.Message.Chat.ID
				return next(ctx, bot, update)
			}

			if _, err := bot.api.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"У вас нет прав.",
			)); err != nil {
				return err
			}
		}

		if update.CallbackQuery != nil {
			if _, ok := adminSet[update.CallbackQuery.From.UserName]; ok {
				return next(ctx, bot, update)
			}

			if _, err := bot.api.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"У вас нет прав.",
			)); err != nil {
				return err
			}
		}

		return nil
	}
}
