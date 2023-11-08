package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
	"strconv"
)

type State interface {
	GetReflectionObject() interface{}
	PerformAdd(message *tgbotapi.Message) (bool, string)
}

type StateRegister struct {
	reflectObj *backend.RegisterRequest
}

func NewStateRegister() *StateRegister {
	return &StateRegister{
		reflectObj: &backend.RegisterRequest{},
	}
}

func (s *StateRegister) GetReflectionObject() interface{} {
	return s.reflectObj
}

func (s *StateRegister) PerformAdd(message *tgbotapi.Message) (bool, string) {
	switch {
	case s.reflectObj.Username == "":
		s.reflectObj.Username = message.Text
		return true, "введите email"
	case s.reflectObj.Email == "":
		s.reflectObj.Email = message.Text
		return true, "введите пароль"
	case s.reflectObj.Password == "":
		s.reflectObj.Password = message.Text
		return false, ""
	default:
		return false, ""
	}
}

type StateLogin struct {
	reflectObj *backend.LoginRequest
}

func NewStateLogin() *StateLogin {
	return &StateLogin{
		reflectObj: &backend.LoginRequest{},
	}
}

func (s *StateLogin) GetReflectionObject() interface{} {
	return s.reflectObj
}

func (s *StateLogin) PerformAdd(message *tgbotapi.Message) (bool, string) {
	switch {
	case s.reflectObj.Username == "":
		s.reflectObj.Username = message.Text
		return true, "введите пароль:"
	case s.reflectObj.Password == "":
		s.reflectObj.Password = message.Text
		return false, ""
	default:
		return false, ""
	}
}

type StateAddChannel struct {
	reflectObj     *backend.ChannelCreateRequest
	PostToChannels *bool
	Access         string
}

func NewStateAddChannel(access string) *StateAddChannel {
	return &StateAddChannel{
		reflectObj: &backend.ChannelCreateRequest{Settings: &backend.ChannelSettingsRequest{}},
		Access:     access,
	}
}

func (s *StateAddChannel) GetReflectionObject() interface{} {
	return s.reflectObj
}

func (s *StateAddChannel) PerformAdd(message *tgbotapi.Message) (bool, string) {
	switch {
	case s.reflectObj.Name == "":
		s.reflectObj.Name = message.Text
		return true, "введите тип канала (тематика). Сейчас поддерживаются 'crypto', 'history', 'sport'"
	case s.reflectObj.Type == "":
		if message.Text != "crypto" && message.Text != "history" && message.Text != "sport" {
			return true, "Вы ввели неподдерживаемую тематику!"
		}
		s.reflectObj.Type = message.Text
		return true, "Введите id вашего канала. Внимание! поменять его потом не получится, введите правильный id"
	case s.reflectObj.OuterID == "":
		s.reflectObj.OuterID = message.Text
		return true, "Приступаем к настройкам. Их можно потом поменять. Введите минимальный рейтинг контента от 1 до 10, который разрешается публиковать в вашем канале"
	case s.reflectObj.Settings.MinRating == 0:
		minRating, err := strconv.Atoi(message.Text)
		if err != nil {
			return true, txtErrMinRating
		}
		if minRating < 1 || minRating > 10 {
			return true, "Значение должно быть в диапазоне от 1 до 10"
		}
		s.reflectObj.Settings.MinRating = minRating
		return true, "Введите число, количество дней спустя которые контент уже не подлежит публикации в ваш канал. По умолчанию 3 дня"
	case s.reflectObj.Settings.NotLaterDays == 0:
		notLaterDays, err := strconv.Atoi(message.Text)
		if err != nil {
			return true, txtErrNotLaterDays
		}
		if notLaterDays == 0 {
			return true, "Значение не может быть равным нулю"
		}
		s.reflectObj.Settings.NotLaterDays = notLaterDays
		return true, "Введите язык для постов. Сейчас поддерживается только 'ru'"
	case s.reflectObj.Settings.Language == "":
		if message.Text != "ru" {
			return true, "Вы ввели неподдерживаемый язык контента"
		}
		s.reflectObj.Settings.Language = message.Text
		return true, txtPostFormatChoice
	case s.PostToChannels == nil:
		txt := "Введите разрешен ли контент без картинок. (да/нет)"
		if message.Text == "0" {
			s.PostToChannels = new(bool)
			*s.PostToChannels = true
			return true, txt
		} else if message.Text == "1" {
			s.PostToChannels = new(bool)
			*s.PostToChannels = false
			return true, txt
		} else {
			return true, "Введите либо 0, либо 1"
		}
	case !s.reflectObj.Settings.EmptyFileAllowed:
		if message.Text != "да" {
			s.reflectObj.Settings.EmptyFileAllowed = true
			return false, ""
		} else if message.Text != "нет" {
			s.reflectObj.Settings.EmptyFileAllowed = false
			return false, ""
		} else {
			return true, "Введите 'да' или 'нет'"
		}
	default:
		return false, ""
	}
}
