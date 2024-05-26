package tg

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/enchik0reo/sup-back/internal/logs"
	"github.com/enchik0reo/sup-back/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Storage interface {
	GetApprovingList(ctx context.Context) ([]models.Approve, error)
	CreateReserved(ctx context.Context, reserve models.Reserved) (int64, error)
	ConfirmApprove(ctx context.Context, id int64, phone string) (int64, error)
	CancelApprove(ctx context.Context, id int64, phone string) (int64, error)

	GetApprovedList(ctx context.Context) ([]models.Approve, error)
	DeleteReserved(ctx context.Context, approveID int64) (int64, error)

	GetPrices(ctx context.Context) ([]models.SupInfo, error)
	EditPrice(ctx context.Context, id, newPrice int64) (int64, error)
	NewSup(ctx context.Context, name string, price int64) (int64, error)
	DeleteSup(ctx context.Context, supID int64) (int64, error)
}

type Bot struct {
	stor Storage

	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
	msgViews map[string]ViewFunc
	admins   map[string]int64
	log      *logs.CustomLog
}

func NewBot(s Storage, token string, admins []string, l *logs.CustomLog) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botApi := Bot{stor: s, api: bot, log: l}

	botApi.addAdmins(admins)

	botApi.createCmdHandlers()

	botApi.createMsgHandlers()

	return &botApi, nil
}

func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	b.log.Info("Telegram bot is getting updates")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		}
	}
}

func (b *Bot) PushNotice() error {
	for _, chatID := range b.admins {
		msg := tgbotapi.NewMessage(chatID, newOrder)

		if _, err := b.api.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			b.log.Error("panic recovered", b.log.Attr("stack", string(debug.Stack())))
		}
	}()

	if update.Message != nil {
		if update.Message.IsCommand() {
			cmd := update.Message.Command()

			cmdView, ok := b.cmdViews[cmd]
			if !ok {
				return
			}

			view := cmdView

			if err := view(ctx, b, update); err != nil {
				b.log.Error("failed to handle update", b.log.Attr("error", err))

				if _, err = b.api.Send(
					tgbotapi.NewMessage(update.Message.Chat.ID, "internal server error"),
				); err != nil {
					b.log.Error("failed to send message", b.log.Attr("error", err))
				}
			}
		} else {
			msg := update.Message.Text

			msgView, ok := b.msgViews[msg]
			if !ok {
				return
			}

			view := msgView

			if err := view(ctx, b, update); err != nil {
				b.log.Error("failed to handle update", b.log.Attr("error", err))

				if _, err = b.api.Send(
					tgbotapi.NewMessage(update.Message.Chat.ID, "internal server error"),
				); err != nil {
					b.log.Error("failed to send message", b.log.Attr("error", err))
				}
			}
		}
	}

	if update.CallbackQuery != nil {

		data := update.CallbackQuery.Data

		dataView, ok := b.msgViews[data]
		if !ok {
			return
		}

		view := dataView

		if err := view(ctx, b, update); err != nil {
			b.log.Error("failed to handle update", b.log.Attr("error", err))

			if _, err = b.api.Send(
				tgbotapi.NewMessage(update.Message.Chat.ID, "internal server error"),
			); err != nil {
				b.log.Error("failed to send message", b.log.Attr("error", err))
			}
		}
	}
}

func (b *Bot) addAdmins(admins []string) {
	b.admins = make(map[string]int64, len(admins))

	for _, a := range admins {
		b.admins[a] = 0
	}
}

func (b *Bot) createCmdHandlers() {
	b.cmdViews = make(map[string]ViewFunc)

	b.addCmdView(startCmd, adminOnly(
		b.admins,
		viewCmdStart()),
	)

	b.addCmdView(showMenuCmd, adminOnly(
		b.admins,
		viewCmdMenu()),
	)

	b.addCmdView(editPriceCmd, adminOnly(
		b.admins,
		viewCmdEditPirce()),
	)

	b.addCmdView(newSupCmd, adminOnly(
		b.admins,
		viewCmdNewSup()),
	)
}

func (b *Bot) createMsgHandlers() {
	b.msgViews = make(map[string]ViewFunc)

	b.addMsgView(reservationList, adminOnly(
		b.admins,
		viewReservationList()),
	)

	b.addMsgView(approvedList, adminOnly(
		b.admins,
		viewApprovedList(),
	))

	b.addMsgView(getPrices, adminOnly(
		b.admins,
		viewPriceList(),
	))
}

func (b *Bot) addCmdView(cmd string, view ViewFunc) {
	b.cmdViews[cmd] = view
}

func (b *Bot) addMsgView(msg string, view ViewFunc) {
	b.msgViews[msg] = view
}

func (b *Bot) deleteMsg(msg string) {
	delete(b.msgViews, msg)
}
