package tg

import (
	"context"
	"fmt"

	"github.com/enchik0reo/sup-back/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func viewPriceList() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		sups, err := bot.stor.GetPrices(ctx)
		if err != nil {
			return err
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(sups))

		for _, sup := range sups {
			info := formatSupInfo(sup)
			data := formatSupData(sup)

			btn := tgbotapi.NewInlineKeyboardButtonData(info, data)

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))

			bot.addMsgView(data, adminOnly(
				bot.admins,
				viewSupOptions(data, sup),
			))
		}

		btn := tgbotapi.NewInlineKeyboardButtonData(addSup, addSup)

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))

		bot.addMsgView(addSup, adminOnly(
			bot.admins,
			viewAddSup(),
		))

		if len(rows) == 0 {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список сапов пуст.")
			msg.ReplyMarkup = keyboard

			if _, err := bot.api.Send(msg); err != nil {
				return err
			}

			return nil
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите существующий или добавьте новый:")
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewSupOptions(data string, sup models.SupInfo) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		data1 := fmt.Sprintf("%s %s", data, editPrice)
		data2 := fmt.Sprintf("%s %s", data, deleteSup)

		btn1 := tgbotapi.NewInlineKeyboardButtonData(editPrice, data1)
		row1 := tgbotapi.NewInlineKeyboardRow(btn1)
		bot.addMsgView(data1, adminOnly(
			bot.admins,
			viewEditPrice(sup.ID, data, data1),
		))

		btn2 := tgbotapi.NewInlineKeyboardButtonData(deleteSup, data2)
		row2 := tgbotapi.NewInlineKeyboardRow(btn2)
		bot.addMsgView(data2, adminOnly(
			bot.admins,
			viewDeleteSup(sup.ID, data, data1, data2),
		))

		keyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2)

		info := formatPriceInfo(sup)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, info)
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewEditPrice(id int64, datas ...string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		defer func() {
			for _, data := range datas {
				bot.deleteMsg(data)
			}
		}()

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
			fmt.Sprintf("Введите команду /%s, id сапа и новую цену:"+
				"\n\nПример:\n/%[1]s %d 1000", editPriceCmd, id))

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewAddSup() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
			fmt.Sprintf("Введите команду /%s, имя сапа и цену:"+
				"\n\nПример:\n/%[1]s BOMBITTO 1000", newSupCmd))

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewDeleteSup(id int64, datas ...string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		defer func() {
			for _, data := range datas {
				bot.deleteMsg(data)
			}
		}()

		if _, err := bot.stor.DeleteSup(ctx, id); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, successDelete)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func formatPriceInfo(sup models.SupInfo) string {
	return fmt.Sprintf("ID: %d\n\nНазвание модели: %s\n\nЦена за сутки в будни: %d₽",
		sup.ID,
		sup.Name,
		sup.Price,
	)
}

func formatSupInfo(sup models.SupInfo) string {
	return fmt.Sprintf("Модель: %s Цена за сутки: %d₽",
		sup.Name,
		sup.Price,
	)
}

func formatSupData(sup models.SupInfo) string {
	return fmt.Sprintf("%d", sup.ID)
}
