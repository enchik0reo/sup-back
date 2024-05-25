package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func viewCmdStart() ViewFunc {
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

		bot.chatID = update.Message.Chat.ID

		keyboard := tgbotapi.NewReplyKeyboard(keys1, keys2, keys3)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меню:")
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
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

		bot.chatID = update.Message.Chat.ID

		keyboard := tgbotapi.NewReplyKeyboard(keys1, keys2, keys3)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меню:")
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
