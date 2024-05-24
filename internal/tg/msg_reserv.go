package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/enchik0reo/sup-back/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type approveData struct {
	models.Approve
	meta string
}

func viewReservationList(storage Storage) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		list, err := storage.GetApproveList(ctx)
		if err != nil {
			return err
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(list))

		for _, elem := range list {
			info := formatInfo(elem)
			data, err := formatData(elem)
			if err != nil {
				return err
			}

			btn := tgbotapi.NewInlineKeyboardButtonData(info, data)

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))

			bot.addMsgView(data, AdminOnly(
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

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список заказов:")
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewReservationOptions(data string, approve models.Approve) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {

		data1, err := addMeta(data, showPhoneNumber)
		if err != nil {
			return err
		}

		data2, err := addMeta(data, approveReserv)
		if err != nil {
			return err
		}

		data3, err := addMeta(data, declineReserv)
		if err != nil {
			return err
		}

		btn1 := tgbotapi.NewInlineKeyboardButtonData(showPhoneNumber, data1)
		row1 := tgbotapi.NewInlineKeyboardRow(btn1)
		bot.addMsgView(data1, AdminOnly(
			bot.admins,
			viewReservationPhone(approve.ClientNumber),
		))

		btn2 := tgbotapi.NewInlineKeyboardButtonData(approveReserv, data2)
		row2 := tgbotapi.NewInlineKeyboardRow(btn2)
		bot.addMsgView(data2, AdminOnly(
			bot.admins,
			viewReservationApprove(approve.ClientNumber, data, data1, data2, data3),
		))

		btn3 := tgbotapi.NewInlineKeyboardButtonData(declineReserv, data3)
		row3 := tgbotapi.NewInlineKeyboardRow(btn3)
		bot.addMsgView(data3, AdminOnly(
			bot.admins,
			viewReservationDecline(approve.ClientNumber),
		))

		keyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Заказ №%d:", approve.ID))
		msg.ReplyMarkup = keyboard

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewReservationPhone(phone string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, phone)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func viewReservationApprove(datas ...string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		defer func() {
			for _, data := range datas {
				bot.deleteMsg(data)
			}
		}()

		ad := approveData{}

		err := json.Unmarshal([]byte(datas[2]), &ad)
		if err != nil {
			return err
		}

		approveToStorage(ctx, bot, ad)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, approved)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func approveToStorage(ctx context.Context, bot *Bot, ad approveData) error {
	_, err := bot.stor.ConfirmApprove(ctx, ad.ID, ad.ClientNumber)
	if err != nil {
		return fmt.Errorf("can't confirm approve: %v", err)
	}

	for _, info := range ad.SupsInfo {
		r := models.Reserved{}
		r.ApproveID = ad.ID
		r.ModelID = info.ID

		temp := info.From
		r.Day = temp

		_, err = bot.stor.CreateReserved(ctx, r)
		if err != nil {
			return fmt.Errorf("can't create reserved: %v", err)
		}

		for temp.Before(info.To) {
			temp = temp.AddDate(0, 0, 1)
			r.Day = temp

			_, err = bot.stor.CreateReserved(ctx, r)
			if err != nil {
				return fmt.Errorf("can't create reserved: %v", err)
			}
		}

		if temp != info.From {
			r.Day = temp

			_, err = bot.stor.CreateReserved(ctx, r)
			if err != nil {
				return fmt.Errorf("can't create reserved: %v", err)
			}
		}
	}

	return nil
}

func viewReservationDecline(datas ...string) ViewFunc {
	return func(ctx context.Context, bot *Bot, update tgbotapi.Update) error {
		defer func() {
			for _, data := range datas {
				bot.deleteMsg(data)
			}
		}()

		ad := approveData{}

		err := json.Unmarshal([]byte(datas[3]), &ad)
		if err != nil {
			return err
		}

		_, err = bot.stor.CancelApprove(ctx, ad.ID, ad.ClientNumber)
		if err != nil {
			return fmt.Errorf("can't confirm approve: %v", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, declined)

		if _, err := bot.api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func formatInfo(approve models.Approve) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("Тел: %s Имя: %s Сумма: %d\n",
		approve.ClientNumber,
		approve.ClientName,
		approve.FullPrice),
	)

	for _, inf := range approve.SupsInfo {
		b.WriteString(fmt.Sprintf("Сап: %s с %s по %s\n",
			inf.Name,
			inf.From.Format(time.DateOnly),
			inf.To.Format(time.DateOnly),
		))
	}

	return b.String()
}

func formatData(approve models.Approve) (string, error) {
	d := approveData{
		approve,
		"",
	}

	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func addMeta(aData, meta string) (string, error) {
	d := approveData{}

	err := json.Unmarshal([]byte(aData), &d)
	if err != nil {
		return "", err
	}

	d.meta = meta

	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
