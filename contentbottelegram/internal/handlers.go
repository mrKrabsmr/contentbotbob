package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go func(update tgbotapi.Update) {
			handling, context := b.Middleware(&update, b.JWTAuthentication)
			if !handling {
				return
			}

			if update.Message == nil {
				return
			}

			b.handleMessage(update.Message, context)
		}(update)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message, context map[string]interface{}) {
	accessToken := context["access_token"].(string)

	if b.Filter(message, b.IsPrivate) {
		if message.IsCommand() {
			switch message.Command() {
			case cmdStart:
				b.handleStart(message)
			}
		} else {
			if b.Filter(message, b.IsNotInState) {
				switch message.Text {
				case cmdRegister:
					b.handleRegister(message)
				case cmdLogin:
					b.handleLogin(message)
				case cmdLogout:
					b.handleLogout(message)
				case cmdProfile:
					b.handleProfile(message, accessToken)
				}
			} else {
				b.handleState(message)
			}
		}
	}
}

func (b *Bot) handleNeedAuth(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, txtNeedAuth)
	msg.ReplyMarkup = authKeyboard

	_, err := b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle NeedAuth", err)
	}
}

func (b *Bot) handleStart(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.ReplyMarkup = mainMenuKeyboard

	_, err := b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Start", err)
	}
}

func (b *Bot) handleState(message *tgbotapi.Message) {
	currState := b.storage[message.From.ID]
	waitNext, txt := currState.PerformAdd(message)
	if waitNext {
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)
		_, err := b.bot.Send(msg)
		if err != nil {
			b.SendErrorMsg(message, "")
			b.logger.Error("***handle State", err)
			return
		}
	} else {
		delete(b.storage, message.From.ID)
		b.handleEndState(currState, message)
	}
}

func (b *Bot) handleEndState(state State, message *tgbotapi.Message) {
	switch state.(type) {

	case *StateRegister:
		result := state.GetReflectionObject()
		if obj, ok := result.(*backend.RegisterRequest); ok {
			_, err := b.api.Register(obj)
			if err != nil {
				b.logger.Error("***handle endRegState", err)
				b.SendErrorMsg(message, err.Error())
				return
			}

			msg := tgbotapi.NewMessage(message.Chat.ID, txtSuccessfulRegister)
			msg.ReplyMarkup = authKeyboard
			_, err = b.bot.Send(msg)
			if err != nil {
				b.SendErrorMsg(message, "")
				b.logger.Error("***handle endRegisterState", err)
				return
			}
		} else {
			b.SendErrorMsg(message, "")
		}

	case *StateLogin:
		result := state.GetReflectionObject()
		if obj, ok := result.(*backend.LoginRequest); ok {
			resp, err := b.api.Login(obj)
			if err != nil {
				b.logger.Error("***handle endLoginState api-login", err)
				b.SendErrorMsg(message, "")
				return
			}

			err = b.redis.SetRefreshToken(message.From.ID, resp.RefreshToken)
			if err != nil {
				b.logger.Error("***handle endLoginState -set-refresh", err)
				b.SendErrorMsg(message, "")
				return
			}

			err = b.redis.SetAccessToken(message.From.ID, resp.AccessToken)
			if err != nil {
				b.logger.Error("***handle endLoginState -set-access", err)
				b.SendErrorMsg(message, "")
				return
			}

			msg := tgbotapi.NewMessage(message.Chat.ID, txtSuccessfulLogin)
			msg.ReplyMarkup = mainMenuKeyboard
			_, err = b.bot.Send(msg)
			if err != nil {
				b.SendErrorMsg(message, "")
				b.logger.Error("***handle endLoginState -send-msg", err)
				return
			}
		} else {
			b.logger.Error("***handle endLoginState", "not convert to obj")
			b.SendErrorMsg(message, "")
			return
		}
	}
}

func (b *Bot) handleRegister(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, txtStartRegister)
	_, err := b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Register", err)
		return
	}
	b.storage[message.From.ID] = NewStateRegister()
}

func (b *Bot) handleLogin(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, txtStartLogin)
	_, err := b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Login", err)
		return
	}
	b.storage[message.From.ID] = NewStateLogin()
}

func (b *Bot) handleLogout(message *tgbotapi.Message) {
	if err := b.redis.DeleteTokens(message.From.ID); err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Logout", err)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, txtSuccessfulLogout)
	msg.ReplyMarkup = authKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Login", err)
		return
	}

}

func (b *Bot) handleBack(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.ReplyMarkup = mainMenuKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle back", err)
		return
	}
}

func (b *Bot) handleProfile(message *tgbotapi.Message, access string) {
	resp, err := b.api.Profile(access)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Profile", err)
		return
	}

	txt := txtProfile(resp)
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	msg.ReplyMarkup = backKeyboard
	_, err = b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Profile", err)
		return
	}

}

//func (b *Bot) handleChannelList(message *tgbotapi.Message) {
//	access, _, err := b.redis.GetTokens(message.From.ID)
//	if err != nil {
//		b.SendErrorMsg(message, "")
//		b.logger.Error("***handle Channel List get access token", err)
//		return
//	}
//
//	resp, err := b.api.ChannelList(access)
//	if err != nil {
//		b.SendErrorMsg(message, "")
//		b.logger.Error("***handle Channel List send response", err)
//		return
//	}
//}
