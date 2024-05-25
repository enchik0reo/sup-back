package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCmd    = "start"
	helpCmd     = "help"
	showMenuCmd = "menu"
)

const (
	newOrder = "📢 Поступил новый заказ!"

	reservationList = "⏳ Заказы на обработку"
	showDetails     = "📋 Детали заказа"
	approveReserv   = "✅ Подтвердить заказ"
	approved        = "👍 Заказ подтвержден"
	declineReserv   = "❌ Отклонить заказ"
	declined        = "👎 Заказ отклонен"

	approvedList     = "🔍 Подтвержденные заказы"
	deleteApproved   = "🗑 Удалить заказ"
	declinedApproved = "🧹 Заказ удален"

	getPrices = "📊 Цены сапов"
	editPrice = "⚙️ Изменить цену"
	DeleteSup = "🏊 Удалить сап"

	addSup = "🏄 Добавить новый сап"
)

type ViewFunc func(ctx context.Context, bot *Bot, update tgbotapi.Update) error
