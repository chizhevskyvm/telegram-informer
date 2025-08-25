package texts

// Сообщения для сценария добавления события
const (
	MsgAskTitle         = "✏️ Введите заголовок события.\nНапример: \"День рождения бабушки 🎉\""
	MsgAskDate          = "📅 Введите дату события в формате ГГГГ-ММ-ДД.\nНапример: \"2025-12-31\""
	MsgAskTime          = "⏰ Введите время события в 24-часовом формате ЧЧ:ММ.\nНапример: \"14:30\""
	MsgConfirm          = "✅ Всё готово! Подтвердите создание события, написав \"да\" или \"нет\"."
	MsgCreated          = "🎉 Событие создано!"
	MsgNoEventsToday    = "📅 Всё чисто! Можно валяться весь день 😎"
	MsgEventsList       = "📅 События на сегодня:"
	MsgDeleteAllError   = "❌ Не удалось отменить события. Попробуйте позже."
	MsgDeleteAllSuccess = "✅ Все события на сегодня отменены."
	MsgDeleteError      = "❌ Не удалось найти событие. Пожалуйста, попробуйте позже."
	MsgDeleteSuccess    = "✅ Событие успешно удалено."
	MsgEventNotFound    = "❌ Не удалось найти событие. Пожалуйста, попробуйте позже."
	MsgError
)

// Кнопки
const (
	BtnMenuTitle            = "Меню:"
	BtnAddEvent             = "➕ Добавить событие"
	BtnTodayEvents          = "📅 Мои события на сегодня"
	BtnAllEvents            = "🗂 Мои события"
	BtnCancelAllTodayEvents = "❌ Отменить все планы на сегодня"
	BtnDeleteEvent          = "🗑 Удалить"
)

// Шаблоны
const (
	BtnEventFormat  = "📌 %s — 🕒 %s" // title, time
	MsgEventDetails = "📅 <b>%s</b>\n🕒 Когда напомнить: <b>%s в %s</b>\n"
)
