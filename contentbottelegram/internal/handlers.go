package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
	"strings"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go func(update tgbotapi.Update) {

			handling, context := b.Middleware(&update, b.JWTAuthentication)
			if !handling {
				return
			}

			if update.CallbackQuery != nil {
				b.handleCallback(update.CallbackQuery, context)
				return
			}

			b.handleMessage(update.Message, context)
		}(update)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message, context map[string]interface{}) {
	obj, needAuth := context["access_token"]

	if b.Filter(message, b.IsPrivate) {
		if message.IsCommand() {
			switch message.Command() {
			case cmdStart:
				b.handleStart(message)
			}
		} else {
			if b.Filter(message, b.IsNotInState) {
				if !needAuth {
					switch message.Text {
					case cmdRegister:
						b.handleRegister(message)
					case cmdLogin:
						b.handleLogin(message)
					}
				} else {
					accessToken := obj.(string)
					switch message.Text {
					case cmdLogout:
						b.handleLogout(message)
					case cmdBack:
						b.handleBack(message)
					case cmdProfile:
						b.handleProfile(message, accessToken)
					case cmdChan:
						b.handleChannelList(message, accessToken)
					}
				}
			} else {
				b.handleState(message)
			}
		}
	}
}

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery, context map[string]interface{}) {
	obj, _ := context["access_token"]
	accessToken := obj.(string)

	if strings.HasPrefix(callback.Data, callbackChannel) {
		b.handleCallbackChannel(callback, accessToken)
	} else if strings.HasPrefix(callback.Data, callbackAddChannel) {
		b.handleCallbackAddChannel(callback, accessToken)
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
	if txt != "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)
		_, err := b.bot.Send(msg)
		if err != nil {
			b.SendErrorMsg(message, "")
			b.logger.Error("***handle State", err)
			return
		}
	}

	if !waitNext {
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
				b.SendErrorMsg(message, err.Error())
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
	case *StateAddChannel:
		stateObj := state.(*StateAddChannel)
		result := stateObj.GetReflectionObject()
		if obj, ok := result.(*backend.ChannelCreateRequest); ok {
			obj.Resource = "telegram"
			err := b.api.ChannelCreate(stateObj.Access, obj)
			if err != nil {
				b.SendErrorMsg(message, err.Error())
				return
			}

			err = b.redis.SetPostFormat(obj.OuterID, *stateObj.PostToChannels)
			if err != nil {
				b.logger.Error("***handle endStateAddChannel", err)
				b.SendErrorMsg(message, "")
				return
			}

			msg := tgbotapi.NewMessage(message.Chat.ID, txtChannelSuccessfulAdded)
			msg.ReplyMarkup = mainMenuKeyboard
			_, err = b.bot.Send(msg)
			if err != nil {
				b.SendErrorMsg(message, "")
				b.logger.Error("***handle endLoginState -send-msg", err)
				return
			}
		} else {
			b.logger.Error("***handle endStateAddChannel", "not convert to obj")
			b.SendErrorMsg(message, "")
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
	msg := tgbotapi.NewMessage(message.Chat.ID, txtBack)
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
		b.SendErrorMsg(message, err.Error())
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

func (b *Bot) handleChannelList(message *tgbotapi.Message, access string) {
	resp, err := b.api.ChannelList(access)
	if err != nil {
		b.SendErrorMsg(message, err.Error())
		return
	}

	txt, keyboardChannels := txtChannels(resp)
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	msg.ReplyMarkup = keyboardChannels
	_, err = b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Channels", err)
		return
	}
}

func (b *Bot) handleCallbackChannel(callback *tgbotapi.CallbackQuery, access string) {
	data, _ := strings.CutPrefix(callback.Data, callbackChannel)

	var newStatus bool
	char := string(data[len(data)-1])

	if char == "1" {
		newStatus = false
	}

	if char == "0" {
		newStatus = true
	}

	channelId, _ := strings.CutSuffix(data, char)

	if err := b.api.ChannelChangeStatus(access, channelId, newStatus); err != nil {
		b.SendErrorMsg(callback.Message, err.Error())
		return
	}

	resp, err := b.api.ChannelList(access)
	if err != nil {
		b.SendErrorMsg(callback.Message, err.Error())
		return
	}

	txt, keyboardChannels := txtChannels(resp)
	msg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, txt)
	msg.ReplyMarkup = &keyboardChannels
	_, err = b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(callback.Message, "")
		b.logger.Error("***handle callbackChannel", err)
		return
	}
}

func (b *Bot) handleCallbackAddChannel(callback *tgbotapi.CallbackQuery, access string) {
	ok, err := b.api.CheckChannelCreate(access)
	if err != nil {
		b.SendErrorMsg(callback.Message, err.Error())
		return
	}

	if !ok {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, txtNotAllowedAddChannel)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.SendErrorMsg(callback.Message, "")
			b.logger.Error("***handle Login", err)
			return
		}

		return
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, txtStartAddChannel)
	_, err = b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(callback.Message, "")
		b.logger.Error("***handle Login", err)
		return
	}

	b.storage[callback.From.ID] = NewStateAddChannel(access)
}
