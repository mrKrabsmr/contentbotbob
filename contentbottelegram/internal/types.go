package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
)

var (
	mainMenuKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdProfile),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdChan),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdSubs),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdLogout),
		),
	)

	authKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdRegister),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdLogin),
		),
	)

	backKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdBack),
		),
	)
)

const (
	cmdStart    = "start"
	cmdRegister = "Регистрация"
	cmdLogin    = "Авторизация"
	cmdLogout   = "Выйти"
	cmdProfile  = "Профиль"
	cmdChan     = "Мои каналы"
	cmdSubs     = "Подписки"
	cmdBack     = "Назад"
)

const (
	txtStartRegister      = "Регистрация. Придумайте уникальный username:"
	txtStartLogin         = "Авторизация. Введите ваш username:"
	txtSuccessfulRegister = "Вы успешно зарегистрировались! На вашу почту был отправлен код подтверждения"
	txtSuccessfulLogin    = "Вы успешно авторизовались!"
	txtSuccessfulLogout   = "Вы вышли из аккаунта"
	txtNeedAuth           = "Для начала необходимо авторизоваться"
)

func txtProfile(profile *backend.ProfileResponse) string {
	txt := fmt.Sprintf(`
Пользователь: %s
Почта: %s
Дата регистрации: %s
`, profile.Username, profile.Email, profile.RegisteredAt)

	if profile.Subscribe != nil {
		sub := profile.Subscribe.(map[string]interface{})
		txt += fmt.Sprintf("Подписка: %s; Осталось: %d", sub["name"].(string), sub["period_until_days"].(int))
	} else {
		txt += "Подписка: нет"
	}

	return txt
}
