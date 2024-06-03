package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
	"strconv"
	"strings"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go func(update tgbotapi.Update) {
			if update.PreCheckoutQuery != nil {
				b.handlePreCheckoutQuery(update.PreCheckoutQuery)
				return
			}

			if update.Message == nil && update.CallbackQuery == nil {
				return
			}

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
					case cmdSendCode:
						b.handleSendCode(message, accessToken)
					case cmdActivate:
						b.handleActivateCode(message, accessToken)
					case cmdProfile:
						b.handleProfile(message, accessToken)
					case cmdChan:
						b.handleChannelList(message, accessToken)
					case cmdSubs:
						b.handleSubscribes(message, accessToken)
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
	} else if strings.HasPrefix(callback.Data, callbackChannelSettings) {
		b.handleCallbackChannelSettingsChange(callback, accessToken)
	} else if strings.HasPrefix(callback.Data, callbackChannelSettingsItem) {
		b.handleCallbackChannelSettingsChangeItem(callback, accessToken)
	}
}

func (b *Bot) handlePreCheckoutQuery(query *tgbotapi.PreCheckoutQuery) {
	g := tgbotapi.PreCheckoutConfig{
		PreCheckoutQueryID: query.ID,
		OK:                 true,
		ErrorMessage:       "",
	}

	_, _ = b.bot.Send(g)
}

func (b *Bot) handleNeedAuth(message *tgbotapi.Message) {
	b.SendMsg(message, txtNeedAuth, authKeyboard)
}

func (b *Bot) handleStart(message *tgbotapi.Message) {
	b.SendMsg(message, "Привет!", mainMenuKeyboard)
}

func (b *Bot) handleState(message *tgbotapi.Message) {
	currState := b.storage[message.From.ID]
	waitNext, txt := currState.PerformAdd(message)
	if txt != "" {
		b.SendMsg(message, txt, nil)
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

			b.SendMsg(message, txtSuccessfulRegister, authKeyboard)
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
			b.SendMsg(message, txtSuccessfulLogin, mainMenuKeyboard)
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

			var chatId int64

			if *stateObj.PostToChannels {
				chanIdConv, err := strconv.Atoi(obj.OuterID)
				if err != nil {
					b.logger.Error("Ошибка при конвертации id канала ", err)
					return
				}

				chanIdConv64 := int64(chanIdConv)
				chatId = chanIdConv64
			} else {
				chatId = message.Chat.ID
			}

			err = b.redis.SetPostChatId(obj.OuterID, chatId)
			if err != nil {
				b.logger.Error("***handle endStateAddChannel", err)
				b.SendErrorMsg(message, "")
				return
			}

			err = b.redis.SetPostPeriod(obj.OuterID, *stateObj.Period)
			if err != nil {
				b.logger.Error("***handle endStateAddChannel", err)
				b.SendErrorMsg(message, "")
				return
			}
			b.SendMsg(message, txtChannelSuccessfulAdded, mainMenuKeyboard)

			go b.StartManage(obj.OuterID, chatId, *stateObj.Period)
		} else {
			b.logger.Error("***handle endStateAddChannel", "not convert to obj")
			b.SendErrorMsg(message, "")
		}

	case *StateActivateCode:
		stateObj := state.(*StateActivateCode)
		result := stateObj.GetReflectionObject()
		if obj, ok := result.(*backend.UserEmailActivateRequest); ok {
			err := b.api.UserEmailActivate(stateObj.Access, obj)
			if err != nil {
				b.SendErrorMsg(message, err.Error())
				return
			}
			b.SendMsg(message, txtSuccessActivate, mainMenuKeyboard)
		} else {
			b.logger.Error("***handle endStateAddChannel", "not convert to obj")
			b.SendErrorMsg(message, "")
		}
	case *StateAPIChannelSettingsChange:
		stateObj := state.(*StateAPIChannelSettingsChange)
		result := stateObj.GetReflectionObject()
		if obj, ok := result.(*backend.ChannelSettingsChangeRequest); ok {
			err := b.api.ChannelSettingsUpdate(stateObj.Access, obj)
			if err != nil {
				b.SendErrorMsg(message, err.Error())
				return
			}

			b.SendMsg(message, txtSuccessSettingsUpdate, mainMenuKeyboard)
		} else {
			b.logger.Error("***handle endStateChannelSettings", "not convert to obj")
			b.SendErrorMsg(message, "")
		}
	case *StateClientChannelSettingsChange:
		stateObj := state.(*StateClientChannelSettingsChange)
		result := stateObj.GetReflectionObject()
		if obj, ok := result.(map[string]interface{}); ok {
			for x, y := range obj {
				if x == "period" {
					err := b.redis.SetPostPeriod(stateObj.ID, y.(int))
					if err != nil {
						b.logger.Error("***handle endStateChannelSettings", err)
						b.SendErrorMsg(message, "")
						return
					}

					chatId, err := b.redis.GetPostChatId(stateObj.ID)
					if err != nil {
						b.logger.Error("***handle endStateChannelSettings", err)
						b.SendErrorMsg(message, "")
						return
					}

					if stop, ok := b.stopManage[stateObj.ID]; ok {
						stop <- struct{}{}
					}

					go b.StartManage(stateObj.ID, chatId, y.(int))
					b.SendMsg(message, txtSuccessSettingsUpdate, mainMenuKeyboard)

				} else if x == "post_to_channels" {
					var chatId int64

					if y.(bool) {
						chanIdConv, err := strconv.Atoi(stateObj.ID)
						if err != nil {
							b.logger.Error("Ошибка при конвертации id канала ", err)
							return
						}

						chanIdConv64 := int64(chanIdConv)
						chatId = chanIdConv64
					} else {
						chatId = message.Chat.ID
					}

					err := b.redis.SetPostChatId(stateObj.ID, chatId)
					if err != nil {
						b.logger.Error("***handle endStateChannelSettings", err)
						b.SendErrorMsg(message, "")
						return
					}

					period, err := b.redis.GetPostPeriod(stateObj.ID)
					if err != nil {
						b.logger.Error("***handle endStateChannelSettings", err)
						b.SendErrorMsg(message, "")
						return
					}

					if stop, ok := b.stopManage[stateObj.ID]; ok {
						stop <- struct{}{}
					}

					go b.StartManage(stateObj.ID, chatId, period)
					b.SendMsg(message, txtSuccessSettingsUpdate, mainMenuKeyboard)
				} else {
					b.logger.Error("***handle endStateChannelSettings", " key invalid")
					b.SendErrorMsg(message, "")
				}
				break
			}
		} else {
			b.logger.Error("***handle endStateChannelSettings", "not convert to obj")
			b.SendErrorMsg(message, "")
		}
	}
}

func (b *Bot) handleRegister(message *tgbotapi.Message) {
	b.SendMsg(message, txtStartRegister, nil)
	b.storage[message.From.ID] = NewStateRegister()
}

func (b *Bot) handleLogin(message *tgbotapi.Message) {
	b.SendMsg(message, txtStartLogin, nil)
	b.storage[message.From.ID] = NewStateLogin()
}

func (b *Bot) handleLogout(message *tgbotapi.Message) {
	if err := b.redis.DeleteTokens(message.From.ID); err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error("***handle Logout", err)
		return
	}

	b.SendMsg(message, txtSuccessfulLogout, authKeyboard)
}

func (b *Bot) handleBack(message *tgbotapi.Message) {
	b.SendMsg(message, txtBack, mainMenuKeyboard)
}

func (b *Bot) handleProfile(message *tgbotapi.Message, access string) {
	resp, err := b.api.Profile(access)
	if err != nil {
		b.SendErrorMsg(message, err.Error())
		return
	}
	txt, isConfirmed := txtProfile(resp)
	if isConfirmed {
		b.SendMsg(message, txt, backKeyboard)
	} else {
		b.SendMsg(message, txt, activateKeyboard)
	}
}

func (b *Bot) handleChannelList(message *tgbotapi.Message, access string) {
	resp, err := b.api.ChannelsUser(access)
	if err != nil {
		b.SendErrorMsg(message, err.Error())
		return
	}

	txt, keyboardChannels := txtChannels(resp)
	b.SendMsg(message, txt, keyboardChannels)

}

func (b *Bot) handleSubscribes(message *tgbotapi.Message, access string) {
	resp, err := b.api.SubscribeList(access)
	if err != nil {
		b.SendErrorMsg(message, err.Error())
		return
	}

	for _, sub := range resp.Result {
		invoice := txtSubscribe(message, sub, b.config.ProviderToken)
		_, err = b.bot.Send(invoice)
		if err != nil {
			b.SendErrorMsg(message, "")
			b.logger.Error(err)
			return
		}
	}
}

func (b *Bot) handleSendCode(message *tgbotapi.Message, access string) {
	if err := b.api.SendActivateCode(access); err != nil {
		b.SendErrorMsg(message, err.Error())
		return
	}

	b.SendMsg(message, txtSuccessSendMail, activateKeyboard)
}

func (b *Bot) handleActivateCode(message *tgbotapi.Message, access string) {
	b.SendMsg(message, txtSetCode, nil)
	b.storage[message.From.ID] = NewStateActivationCode(access)
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

	if newStatus {
		chatId, err := b.redis.GetPostChatId(channelId)
		if err != nil {
			b.logger.Error(err)
			b.SendErrorMsg(callback.Message, "")
			return
		}

		postPeriod, err := b.redis.GetPostPeriod(channelId)
		if err != nil {
			b.logger.Error(err)
			b.SendErrorMsg(callback.Message, "")
			return
		}
		go b.StartManage(channelId, chatId, postPeriod)

	} else {
		if stop, ok := b.stopManage[channelId]; ok {
			stop <- struct{}{}
		}
	}

	if err := b.api.ChannelChangeStatus(access, channelId, newStatus); err != nil {
		b.SendErrorMsg(callback.Message, err.Error())
		return
	}

	resp, err := b.api.ChannelsUser(access)
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
		b.logger.Error(err)
	}
}

func (b *Bot) handleCallbackAddChannel(callback *tgbotapi.CallbackQuery, access string) {
	ok, err := b.api.CheckChannelCreate(access)
	if err != nil {
		b.SendErrorMsg(callback.Message, err.Error())
		return
	}

	if !ok {
		b.SendMsg(callback.Message, txtNotAllowedAddChannel, nil)
		return
	}

	b.SendMsg(callback.Message, txtStartAddChannel, nil)
	b.storage[callback.From.ID] = NewStateAddChannel(access)
}

func (b *Bot) handleCallbackChannelSettingsChange(callback *tgbotapi.CallbackQuery, access string) {
	data, _ := strings.CutPrefix(callback.Data, callbackChannelSettings)

	settings, err := b.api.ChannelSettings(access, data)
	if err != nil {
		b.SendErrorMsg(callback.Message, "")
		b.logger.Error(err)
		return
	}

	postPeriod, err := b.redis.GetPostPeriod(data)
	if err != nil {
		b.SendErrorMsg(callback.Message, "")
		b.logger.Error(err)
		return
	}

	chatId, err := b.redis.GetPostChatId(data)
	if err != nil {
		b.SendErrorMsg(callback.Message, "")
		b.logger.Error(err)
		return
	}

	var postToChan bool

	if chatId == callback.Message.Chat.ID {
		postToChan = false
	} else {
		postToChan = true
	}

	if _, is := someFunc(postPeriod, b.config.SecretNum); is {
		postPeriod = -11111
	}

	txt, keyboard := txtChannelSettingsChange(data, settings, postPeriod, postToChan)
	b.SendMsg(callback.Message, txt, keyboard)
}

func (b *Bot) handleCallbackChannelSettingsChangeItem(callback *tgbotapi.CallbackQuery, access string) {
	data, _ := strings.CutPrefix(callback.Data, callbackChannelSettingsItem)
	items := strings.Split(data, ":-")
	chanId := items[0]
	changeField := items[1]

	if changeField == "min_rating" ||
		changeField == "not_later_days" ||
		changeField == "empty_file_allowed" ||
		changeField == "language" {

		state := NewStateAPIChannelSettingsChange(access)
		state.reflectObj.ID = chanId
		state.reflectObj.Data[changeField] = nil
		b.storage[callback.From.ID] = state

		var txt string
		switch changeField {
		case "min_rating":
			txt = strings.TrimPrefix(txtSetMinRating, "Приступаем к настройкам. Их можно потом поменять. ")
		case "not_later_days":
			txt = txtSetNotLaterDays
		case "empty_file_allowed":
			txt = txtSetAllowedFileEmpty
		case "language":
			txt = txtSetLanguage
		}

		b.SendMsg(callback.Message, txt, nil)
	} else if changeField == "period" ||
		changeField == "post_to_channels" {

		state := NewStateClientAPIChannelSettingsChange(access)
		state.ID = chanId
		state.reflectObj[changeField] = nil
		b.storage[callback.From.ID] = state

		var txt string
		switch changeField {
		case "period":
			txt = txtSetPeriod
		case "post_to_channels":
			txt = txtSetPostToChannels
		}

		b.SendMsg(callback.Message, txt, nil)
	} else {
		b.SendErrorMsg(callback.Message, "")
		b.logger.Error("***handle callback channel settings item")
		return
	}
}
