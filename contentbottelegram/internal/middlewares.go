package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) JWTAuthentication(update *tgbotapi.Update, context *map[string]interface{}) bool {
	if b.Filter(update.Message, b.IsPrivate, b.IsNotAuth, b.IsNotInState) {
		userTgId := update.Message.From.ID
		access, refresh, err := b.redis.GetTokens(userTgId)
		if err != nil {
			b.logger.Error("***middleware JWTAuthentication (get token)", err)
			return false
		}

		if refresh == "" {
			b.handleNeedAuth(update.Message)
			return false
		}

		if access == "" {
			resp, err := b.api.UpdateAccessToken(refresh)
			if err != nil {
				b.logger.Error("***middleware JWTAuthentication (update access token)", err)
				return false
			}

			if err = b.redis.SetAccessToken(userTgId, resp.AccessToken); err != nil {
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
