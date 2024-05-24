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
	msgHelp = `/menu - команда для показа меню`

	NewOrder = "📢 Поступил новый заказ!"

	reservationList = "⏳ Заказы на подтверждение" // Пункт в меню
	// Когда тыкаешь на ReservationList, показывается номер и вылазит 3 опции:
	showPhoneNumber = "☎️ Номер телефона"   // Подпункт на экране (для каждой записи)
	approveReserv   = "✅ Подтвердить заказ" // Подпункт на экране (для каждой записи)
	approved        = "👍 Заказ подтвержден"
	declineReserv   = "❌ Отклонить заказ" // Подпункт на экране  (для каждой записи)
	declined        = "👎 Заказ отклонен"

	approvedList = "🔍 Подтвержденные заказы" // Пункт в меню
	// Когда тыкаешь на ApprovedList, показывается номер и вылазит 1 опция:
	DeleteApprovedCmd = "🗑 Удалить действующий заказ" // Подпункт на экране

	getPrices = "📊 Цены сапов" // Пункт в меню
	// Когда тыкаешь на GetPrices вылазит 1 опция:
	EditPricesCmd = "⚙️ Изменить цену" // Подпункт на экране
)

type ViewFunc func(ctx context.Context, bot *Bot, update tgbotapi.Update) error
