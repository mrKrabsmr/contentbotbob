package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
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

func (b *Bot) StartManage(chanId string, chatId int64, sleepTime int) {
	stop := make(chan struct{})
	b.stopManage[chanId] = stop

	defer delete(b.stopManage, chanId)

	duration := time.Minute * time.Duration(sleepTime)
	if num, is := someFunc(sleepTime, b.config.SecretNum); is {
		duration = time.Second * time.Duration(num)
	}

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-stop:
			break LOOP
		case <-ticker.C:
			resp, ok, err := b.api.Content(chanId)
			if err != nil {
				b.logger.Error("Ошибка при попытке ведения канала ", err)
			} else if !ok {
				break LOOP
			} else {
				if len(resp.Result.Images) > 0 {
					var mediaGroup []interface{}

					wasAddedText := false
					for _, image := range resp.Result.Images {
						inputMediaPhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(image.FileUrl))
						if !wasAddedText {
							inputMediaPhoto.Caption = resp.Result.Text
							wasAddedText = true
						}

						mediaGroup = append(mediaGroup, inputMediaPhoto)
					}

					msgGroup := tgbotapi.NewMediaGroup(chatId, mediaGroup)
					_, _ = b.bot.SendMediaGroup(msgGroup)
				} else {
					msgTxt := tgbotapi.NewMessage(chatId, resp.Result.Text)
					_, _ = b.bot.Send(msgTxt)
				}
			}
		}
	}
}

func (b *Bot) SendMsg(message *tgbotapi.Message, txt string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	_, err := b.bot.Send(msg)
	if err != nil {
		b.SendErrorMsg(message, "")
		b.logger.Error(err)
		return
	}
}

func (b *Bot) SendErrorMsg(message *tgbotapi.Message, txt string) {
	if txt == "" {
		txt = txtBotError
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	_, _ = b.bot.Send(msg)
}

func someFunc(i int, a int) (int, bool) {
	val := i - a
	if val%100000 == 0 {
		return val / 100000, true
	}

	return -1, false
}
