package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) JWTAuthentication(update *tgbotapi.Update, context *map[string]interface{}) bool {
	var message *tgbotapi.Message
	var userTgID int64

	if update.Message != nil {
		message = update.Message
		userTgID = update.Message.From.ID
	} else if update.CallbackQuery != nil {
		message = update.CallbackQuery.Message
		userTgID = update.CallbackQuery.From.ID
	}

	if b.Filter(message, b.IsPrivate, b.IsNotAuth, b.IsNotInState) {
		access, refresh, err := b.redis.GetTokens(userTgID)
		if err != nil {
			b.logger.Error("***middleware JWTAuthentication (get token)", err)
			return false
		}

		if refresh == "" {
			b.handleNeedAuth(message)
			return false
		}

		if access == "" {
			resp, err := b.api.UpdateAccessToken(refresh)
			if err != nil {
				b.logger.Error("***middleware JWTAuthentication (update access token)", err)
				return false
			}

			if err = b.redis.SetAccessToken(userTgID, resp.AccessToken); err != nil {
				b.logger.Error("***middleware JWTAuthentication (set access token)", err)
				return false
			}

			(*context)["access_token"] = resp.AccessToken
		} else {
			(*context)["access_token"] = access
		}

		return true
	}

	return true
}
