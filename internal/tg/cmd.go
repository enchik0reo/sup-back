package tg

import (
	"context"
	"strconv"
	"strings"

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

		keyboard := tgbotapi.NewReplyKeyboard(keys1, keys2, keys3)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меню:")
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewCmdEditPirce() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		data := update.Message.CommandArguments()

		datas := strings.Split(data, " ")

		id, err := strconv.Atoi(datas[0])
		if err != nil {
			return err
		}

		price, err := strconv.Atoi(datas[1])
		if err != nil {
			return err
		}

		if _, err := bot.stor.EditPrice(ctx, int64(id), int64(price)); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, successEditPrice)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewCmdNewSup() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		data := update.Message.CommandArguments()

		datas := strings.Split(data, " ")

		price, err := strconv.Atoi(datas[1])
		if err != nil {
			return err
		}

		if _, err := bot.stor.NewSup(ctx, datas[0], int64(price)); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, successAdd)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
