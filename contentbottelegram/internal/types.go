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
	cmdAddChan  = "Добавить канал"
)

const (
	txtStartRegister      = "Регистрация. Придумайте уникальный username:"
	txtStartLogin         = "Авторизация. Введите ваш username:"
	txtSuccessfulRegister = `Вы успешно зарегистрировались! На вашу почту был отправлен код подтверждения.
Авторизуйтесь и в разделе "Профиль" подвердите свой email.`
	txtSuccessfulLogin  = "Вы успешно авторизовались!"
	txtSuccessfulLogout = "Вы вышли из аккаунта"
	txtNeedAuth         = "Для начала необходимо авторизоваться"
	txtBack             = "Принято. На главное меню"
	txtErrMinRating     = "Введите, пожалуйста корректное число от 1 до 10"
	txtErrNotLaterDays  = "Введите, пожалуйста корректное число"
	txtStartAddChannel  = `Чтобы добавить канал вы должны указывать подлинные данные последовательно. 
Отвечайте на сообщения бота четко, как положено.

Обратите внимание, что название, тематика и платформа канала, а также его id в этой платформе не могут в дальнейшем быть изменены.
ПРОСЬБА ОТНЕСТИТЬ К ЭТОМУ ВНИМАТЕЛЬНО!

Введите название канала:`
	txtPostFormatChoice = `Выберите вариант, который вам больше подходит (напишите номер варианта)
0. Бот будет вести непосредственно ваш канал самолично
1. Бот будет отправлять посты вам личным сообщением`
	txtNotAllowedAddChannel   = "У вас либо нет подписки, либо вы уже добавили максимально возможное количество каналов по вашей подписке"
	txtChannelsEmpty          = "Вами еще не был добавлен ни один канал"
	txtChannelSuccessfulAdded = "Канал успешно добавлен"
)

const (
	callbackChannel    = "channel"
	callbackAddChannel = "addChan"
)

func txtProfile(profile *backend.ProfileResponse) string {
	txt := fmt.Sprintf(`
Пользователь: %s
Почта: %s
Дата регистрации: %s
`, profile.Username, profile.Email, profile.RegisteredAt)

	if profile.Subscribe != nil {
		sub := profile.Subscribe
		txt += fmt.Sprintf(
			`Подписка: %s;
																	 	 максимум каналов: %d`,
			sub["name"].(string), int(sub["max_use_channels"].(float64)),
		)
	} else {
		txt += "Подписка: нет"
	}

	return txt
}

func txtChannels(channels []*backend.ChannelResponse) (string, tgbotapi.InlineKeyboardMarkup) {
	var buttons [][]tgbotapi.InlineKeyboardButton
	var txt string
	for num, channel := range channels {
		var statusString, statusTxt, toDoTxt string
		if channel.StatusOn {
			statusString = "1"
			statusTxt = "🟢"
			toDoTxt = "выключить"
		} else {
			statusString = "0"
			statusTxt = "🔴"
			toDoTxt = "включить"
		}

		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(channel.Name+" | "+toDoTxt, callbackChannel+channel.ID+statusString)),
		)

		txt += fmt.Sprintf(`
%d. %s
Тип: %s
Платформа: %s
Статус: %s
`,
			num+1, channel.Name, channel.Type, channel.Resource, statusTxt,
		)
	}

	addChanButton := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(cmdAddChan, callbackAddChannel))
	buttons = append(buttons, addChanButton)

	if txt == "" {
		txt = txtChannelsEmpty
	}

	return txt, tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
