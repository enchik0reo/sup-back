package tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func viewCmdHelp() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		if _, err := bot.api.Send(tgbotapi.NewMessage(
			update.FromChat().ID,
			msgHelp)); err != nil {
			return err
		}

		return nil
	}
}

func viewCmdStart() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		if _, err := bot.api.Send(tgbotapi.NewMessage(
			update.FromChat().ID,
			fmt.Sprintf("Hello, %s!", update.FromChat().UserName))); err != nil {
			return err
		}

		return nil
	}
}

func viewCmdMenu() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		keys1 := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(reservationList),
		}

		keys2 := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(approvedList),
		}

		keys3 := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(getPrices),
		}

		keyboard := tgbotapi.NewReplyKeyboard(keys1, keys2, keys3)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меню:")
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
