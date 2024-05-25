package tg

import (
	"context"
	"fmt"

	"github.com/enchik0reo/sup-back/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func viewPriceList(storage Storage) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		sups, err := storage.GetPrices(ctx)
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

		// Для добавления нового сапа будет строчка в конце списка сапов !!!

		// info := addSup
		// data := addSup

		// btn := tgbotapi.NewInlineKeyboardButtonData(info, data)

		// rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))

		// bot.addMsgView(data, adminOnly(
		// 	bot.admins,
		// 	viewAddSup(data, sup),
		// ))

		if len(rows) == 0 {
			// keyboard := tgbotapi.NewInlineKeyboardMarkup(btn)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список сапов пуст.")
			// msg.ReplyMarkup = keyboard

			if _, err := bot.api.Send(msg); err != nil {
				return err
			}

			return nil
		}

		// rows = append(rows, btn)

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите сап или добавьте новый:")
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

		btn1 := tgbotapi.NewInlineKeyboardButtonData(deleteApproved, data1)
		row1 := tgbotapi.NewInlineKeyboardRow(btn1)
		bot.addMsgView(data1, adminOnly(
			bot.admins,
			viewEditPrice(sup, data, data1),
		))

		// TODO Добавить кнопку удаления сапа viewDeleteSup

		keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)

		info := formatPriceInfo(sup)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, info)
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

// TODO ...
func viewEditPrice(sup models.SupInfo, datas ...string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		defer func() {
			for _, data := range datas {
				bot.deleteMsg(data)
			}
		}()

		_, err := bot.stor.DeleteReserved(ctx, approve.ID)
		if err != nil {
			return fmt.Errorf("can't delete reserved: %v", err)
		}

		_, err = bot.stor.CancelApprove(ctx, approve.ID, approve.ClientNumber)
		if err != nil {
			return fmt.Errorf("can't delete reserved: %v", err)
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, declinedApproved)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

// TODO viewAddSup

// TODO viewDeleteSup

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
