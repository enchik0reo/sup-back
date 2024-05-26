package tg

import (
	"context"
	"fmt"

	"github.com/enchik0reo/sup-back/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func viewApprovedList() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		list, err := bot.stor.GetApprovedList(ctx)
		if err != nil {
			return err
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(list))

		for _, elem := range list {
			info := formatOrderInfo(elem)
			data := formatData(elem)

			btn := tgbotapi.NewInlineKeyboardButtonData(info, data)

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))

			bot.addMsgView(data, adminOnly(
				bot.admins,
				viewApproveOptions(data, elem),
			))
		}

		if len(rows) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Подтвержденных заказов нет.")

			if _, err := bot.api.Send(msg); err != nil {
				return err
			}

			return nil
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Последние подтвержденные заказы:")
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewApproveOptions(data string, approve models.Approve) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		data1 := fmt.Sprintf("%s %s", data, deleteApproved)

		btn1 := tgbotapi.NewInlineKeyboardButtonData(deleteApproved, data1)
		row1 := tgbotapi.NewInlineKeyboardRow(btn1)
		bot.addMsgView(data1, adminOnly(
			bot.admins,
			viewApproveDelete(approve, data, data1),
		))

		keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)

		info := formatInfo(approve)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, info)
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewApproveDelete(approve models.Approve, datas ...string) ViewFunc {
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
