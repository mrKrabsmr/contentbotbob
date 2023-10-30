package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) IsPrivate(message *tgbotapi.Message) bool {
	if message.Chat.Type == "private" {
		return true
	}

	return false
}

func (b *Bot) IsNotInState(message *tgbotapi.Message) bool {
	if _, ok := b.storage[message.From.ID]; !ok {
		return true
	}

	return false
}

func (b *Bot) IsNotAuth(message *tgbotapi.Message) bool {
	if message.Text == cmdRegister || message.Text == cmdLogin {
		return false
	}

	return true
}
