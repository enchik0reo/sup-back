package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCmd     = "start"
	showMenuCmd  = "menu"
	editPriceCmd = "price"
	newSupCmd    = "new"
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

	getPrices        = "📝 Список сапов"
	editPrice        = "📊 Изменить цену"
	successEditPrice = "💵 Цена изменена"
	deleteSup        = "🏊 Удалить сап"
	successDelete    = "🧹 Сап удален"

	addSup     = "🏄 Добавить новый сап"
	successAdd = "🎉 Сап добавлен"
)

type ViewFunc func(ctx context.Context, bot *Bot, update tgbotapi.Update) error
