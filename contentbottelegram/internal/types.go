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

	activateKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdSendCode),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(cmdActivate),
		),
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
	cmdSendCode = "Отправить код на почту"
	cmdActivate = "Ввести код подверждаения"
)

const (
	txtSetEmail           = "Введите email:"
	txtSetPassword        = "Введите пароль:"
	txtStartRegister      = "Регистрация. Придумайте уникальный username:"
	txtStartLogin         = "Авторизация. Введите ваш username:"
	txtSuccessfulRegister = `Вы успешно зарегистрировались! На вашу почту был отправлен код подтверждения.
Авторизуйтесь и в разделе "Профиль" подвердите свой email.`
	txtSuccessfulLogin       = "Вы успешно авторизовались!"
	txtSuccessfulLogout      = "Вы вышли из аккаунта"
	txtNeedAuth              = "Для начала необходимо авторизоваться"
	txtBack                  = "Принято. На главное меню"
	txtSuccessActivate       = "Аккаунт успешно подвержден, вы можете пользоваться всеми функциями бота"
	txtSuccessSettingsUpdate = "Настройки успешно применены"
	txtSuccessSendMail       = "сообщение успешно отправлено вам на почту"
	txtSetCode               = "Введите код: "
	txtStartAddChannel       = `Чтобы добавить канал вы должны указывать подлинные данные последовательно. 
Отвечайте на сообщения бота четко, как положено.

Обратите внимание, что название, тематика канала, а также его id не могут в дальнейшем быть изменены.
ПРОСЬБА ОТНЕСТИТЬ К ЭТОМУ ВНИМАТЕЛЬНО!

Введите название канала:`
	txtUnSupportType       = "Вы ввели неподдерживаемую тематику!"
	txtSetId               = "Введите id вашего канала. Внимание! поменять его потом не получится, введите правильный id"
	txtIncorrectId         = "Неккоректное значение id"
	txtSetMinRating        = "Приступаем к настройкам. Их можно потом поменять. Введите минимальный рейтинг контента от 1 до 10, который разрешается публиковать в вашем канале"
	txtErrMinRating        = "Введите, пожалуйста корректное число от 1 до 10"
	txtErrMinRating2       = "Значение должно быть в диапазоне от 1 до 10"
	txtSetNotLaterDays     = "Введите число, количество дней спустя которые контент уже не подлежит публикации в ваш канал. По умолчанию 3 дня"
	txtErrNotLaterDays     = "Введите, пожалуйста корректное число"
	txtErrNotLaterDays2    = "Значение не может быть равным нулю"
	txtSetLanguage         = "Введите язык для постов. Сейчас поддерживается только 'ru'"
	txtErrLanguage         = "Вы ввели неподдерживаемый язык контента"
	txtSetPeriod           = "Введите периодичность отправления постов в минутах. Не менее десяти минут!"
	txtErrPeriod           = "Периодичность - это числовое значение"
	txtErrPeriod2          = "Периодичность не может быть меньше десяти минут"
	txtSetPostFormatChoice = `Выберите вариант, который вам больше подходит (напишите номер варианта)
0. Бот будет вести непосредственно ваш канал самолично
1. Бот будет отправлять посты вам личным сообщением`
	txtSetPostToChannels      = "Введите, отправлять контент в канал? (да/нет) (при 'нет' контент будет приходить вам в личные сообщения) "
	txtSetAllowedFileEmpty    = "Введите разрешен ли контент без картинок. (да/нет)"
	txtErrPostToChannels      = "Введите либо 0, либо 1"
	txtErrAllowedFileEmpty    = "Введите 'да' или 'нет'"
	txtNotAllowedAddChannel   = "У вас либо нет подписки, либо вы уже добавили максимально возможное количество каналов по вашей подписке"
	txtChannelsEmpty          = "Вами еще не был добавлен ни один канал"
	txtChannelSuccessfulAdded = "Канал успешно добавлен"
	txtBotError               = "Извините, произошла ошибка в работе бота"
)

const (
	callbackChannel             = "channel"
	callbackChannelSettings     = "settingsChannel"
	callbackChannelSettingsItem = "itemChannelSettings"
	callbackAddChannel          = "addChan"
	callbackSub                 = "subscribe"
)

var useChanTypes = []string{
	"all",
	"crypto",
	"films",
	"animals",
    "retails",
	"cybersport",
	"sport",
	"football",
	"basketball",
	"volleyball",
	"tennis",
	"boxing",
	"auto",
	"figureskating",
	"hockey",
}

func txtProfile(profile *backend.ProfileResponse) (string, bool) {
	var confirmTxt string
	if profile.IsConfirmed {
		confirmTxt = "✅"
	} else {
		confirmTxt = "❌"
	}

	txt := fmt.Sprintf(`
Пользователь: %s
Почта: %s
Дата регистрации: %s
Аккаунт подтвержден: %s
`, profile.Username, profile.Email, profile.RegisteredAt, confirmTxt)

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

	if !profile.IsConfirmed {

	}

	if !profile.IsConfirmed {
		return txt, false
	}

	return txt, true
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
			tgbotapi.NewInlineKeyboardButtonData(
				channel.Name+" | "+toDoTxt,
				callbackChannel+channel.OuterID+statusString),
		),
		)

		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf(
				"Изменить настройки «%s»", channel.Name),
				callbackChannelSettings+channel.OuterID),
		),
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

func txtChannelSettingsChange(chanId string, settings *backend.ChannelSettingsResponse, period int, postToChannels bool) (string, tgbotapi.InlineKeyboardMarkup) {
	var allowedEmptyFile string
	if settings.EmptyFileAllowed {
		allowedEmptyFile = "Разрешено"
	} else {
		allowedEmptyFile = "Запрещено"
	}

	var postTo string
	if postToChannels {
		postTo = "да"
	} else {
		postTo = "нет"
	}

	items := map[string]interface{}{
		"min_rating":         fmt.Sprintf("Мин. рейтинг (%d)", settings.MinRating),
		"empty_file_allowed": fmt.Sprintf("Разр. без медиа-файлов (%s)", allowedEmptyFile),
		"not_later_days":     fmt.Sprintf("Срок годн. (%d д.)", settings.NotLaterDays),
		"language":           fmt.Sprintf("Язык контента (%s)", settings.Language),
		"period":             fmt.Sprintf("Периодичность (раз в %d м.)", period),
		"post_to_channels":   fmt.Sprintf("Сразу в канал (%s)", postTo),
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for key, item := range items {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(item.(string), callbackChannelSettingsItem+chanId+":-"+key)),
		)
	}

	return "Выберите, что хотите поменять", tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func txtSubscribe(message *tgbotapi.Message, sub *backend.SubscribeResponse, provToken string) tgbotapi.InvoiceConfig {
	val, _ := sub.PriceRub.Float64()

	title := fmt.Sprintf("%s (макс. каналов: %d)", sub.Name, sub.MaxUseChannels)

	label := []tgbotapi.LabeledPrice{tgbotapi.LabeledPrice{
		Label:  title,
		Amount: int(val) * 100,
	}}

	invoice := tgbotapi.NewInvoice(
		message.Chat.ID,
		title,
		sub.Description,
		"test-payload",
		provToken,
		"kekekek",
		"rub",
		label,
	)
	invoice.MaxTipAmount = 50 * 100
	invoice.SuggestedTipAmounts = []int{5 * 100, 25 * 100, 50 * 100}

	return invoice
}

func txtChannelTypesAllowed() string {
	txt := "введите тип канала (тематика). Сейчас поддерживаются:\n"
	for _, element := range useChanTypes {
		txt += fmt.Sprintf("'%s'\n", element)
	}

	return txt
}
