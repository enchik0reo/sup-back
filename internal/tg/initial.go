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
	GetApproveList(ctx context.Context) ([]models.Approve, error)
	CreateReserved(ctx context.Context, reserve models.Reserved) (int64, error)
	ConfirmApprove(ctx context.Context, id int64, phone string) (int64, error)
	CancelApprove(ctx context.Context, id int64, phone string) (int64, error)
}

type Bot struct {
	stor Storage

	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
	msgViews map[string]ViewFunc
	admins   map[string]struct{}
	log      *logs.CustomLog
}

func NewBot(stor Storage, token string, admins []string, l *logs.CustomLog) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botApi := Bot{api: bot, stor: stor, log: l}

	botApi.addAdmins(admins)

	botApi.createBasicCmdMenu()

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

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			b.log.Error("panic recovered", b.log.Attr("recover", p), b.log.Attr("stack", debug.Stack()))
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
	b.admins = make(map[string]struct{}, len(admins))

	for _, a := range admins {
		b.admins[a] = struct{}{}
	}
}

func (b *Bot) createBasicCmdMenu() {
	b.cmdViews = make(map[string]ViewFunc)
	b.msgViews = make(map[string]ViewFunc)

	b.addCmdView(startCmd, AdminOnly(
		b.admins,
		viewCmdStart()),
	)

	b.addCmdView(helpCmd, AdminOnly(
		b.admins,
		viewCmdHelp()),
	)

	b.addCmdView(showMenuCmd, AdminOnly(
		b.admins,
		viewCmdMenu()),
	)

	b.addMsgView(reservationList, AdminOnly(
		b.admins,
		viewReservationList(b.stor)),
	)

	// TODO approvedList

	// TODO getPrices
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
