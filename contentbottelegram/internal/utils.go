package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Middleware(update *tgbotapi.Update, middlewares ...func(*tgbotapi.Update, *map[string]interface{}) bool) (bool, map[string]interface{}) {
	context := make(map[string]interface{})
	for _, middleware := range middlewares {
		isOk := middleware(update, &context)
		if !isOk {
			return false, nil
		}
	}

	return true, context
}

func (b *Bot) Filter(message *tgbotapi.Message, filters ...func(message *tgbotapi.Message) bool) bool {
	for _, filter := range filters {
		isOk := filter(message)
		if !isOk {
			return false
		}
	}

	return true
}

func (b *Bot) SendErrorMsg(message *tgbotapi.Message, txt string) {
	if txt == "" {
		txt = "Извините, произошла ошибка в работе бота"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	_, _ = b.bot.Send(msg)
}
