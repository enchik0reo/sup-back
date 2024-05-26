package tg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/enchik0reo/sup-back/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func viewReservationList() ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		list, err := bot.stor.GetApprovingList(ctx)
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
				viewReservationOptions(data, elem),
			))
		}

		if len(rows) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Новых заказов нет.")

			if _, err := bot.api.Send(msg); err != nil {
				return err
			}

			return nil
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Необработанных заказов: %d", len(rows)))
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewReservationOptions(data string, approve models.Approve) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		data2 := fmt.Sprintf("%s %s", data, approveReserv)
		data3 := fmt.Sprintf("%s %s", data, declineReserv)

		btn2 := tgbotapi.NewInlineKeyboardButtonData(approveReserv, data2)
		row2 := tgbotapi.NewInlineKeyboardRow(btn2)
		bot.addMsgView(data2, adminOnly(
			bot.admins,
			viewReservationApprove(approve, data, data2, data3),
		))

		btn3 := tgbotapi.NewInlineKeyboardButtonData(declineReserv, data3)
		row3 := tgbotapi.NewInlineKeyboardRow(btn3)
		bot.addMsgView(data3, adminOnly(
			bot.admins,
			viewReservationDecline(approve, data, data2, data3),
		))

		keyboard := tgbotapi.NewInlineKeyboardMarkup(row2, row3)

		info := formatInfo(approve)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, info)
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewReservationApprove(approve models.Approve, datas ...string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		defer func() {
			for _, data := range datas {
				bot.deleteMsg(data)
			}
		}()

		approveToStorage(ctx, bot, approve)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, approved)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func approveToStorage(ctx context.Context, bot *Bot, approve models.Approve) error {
	_, err := bot.stor.ConfirmApprove(ctx, approve.ID, approve.ClientNumber)
	if err != nil {
		return fmt.Errorf("can't confirm approve: %v", err)
	}

	for _, info := range approve.SupsInfo {
		r := models.Reserved{}
		r.ApproveID = approve.ID
		r.ModelID = info.ID

		temp := info.From
		r.Day = temp

		_, err = bot.stor.CreateReserved(ctx, r)
		if err != nil {
			return fmt.Errorf("can't create reserved: %v", err)
		}

		for {
			temp = temp.AddDate(0, 0, 1)
			r.Day = temp
			if temp.Before(info.To) {
				_, err = bot.stor.CreateReserved(ctx, r)
				if err != nil {
					return fmt.Errorf("can't create reserved: %v", err)
				}
			} else {
				break
			}
		}
	}

	return nil
}

func viewReservationDecline(approve models.Approve, datas ...string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		defer func() {
			for _, data := range datas {
				bot.deleteMsg(data)
			}
		}()

		_, err := bot.stor.CancelApprove(ctx, approve.ID, approve.ClientNumber)
		if err != nil {
			return fmt.Errorf("can't decline approve: %v", err)
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, declined)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func formatOrderInfo(approve models.Approve) string {
	return fmt.Sprintf("Заказ №%d от %s на сумму %d₽",
		approve.ID,
		approve.ClientName,
		approve.FullPrice,
	)
}

func formatInfo(approve models.Approve) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("Заказ №%d\n\nИмя: %s\n\nТел: %s\n\n",
		approve.ID,
		approve.ClientName,
		approve.ClientNumber,
	))

	b.WriteString("Бронирование:\n\n")

	for _, inf := range approve.SupsInfo {
		b.WriteString(fmt.Sprintf("%s:\nс %s по %s\n\n",
			inf.Name,
			inf.From.Format(time.DateOnly),
			inf.To.Format(time.DateOnly),
		))
	}

	b.WriteString(fmt.Sprintf("Сумма: %d₽\n\n", approve.FullPrice))

	return b.String()
}

func formatData(approve models.Approve) string {
	return fmt.Sprintf("%d %s", approve.ID, approve.ClientNumber)
}
