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
		return true, txtSetEmail
	case s.reflectObj.Email == "":
		s.reflectObj.Email = message.Text
		return true, txtSetPassword
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
		return true, txtSetPassword
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
	Period         *int
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
		return true, txtChannelTypesAllowed()
	case s.reflectObj.Type == "":
		ok := false
		for _, element := range useChanTypes {
			if element == message.Text {
				ok = true
			}
		}

		if !ok {
			return true, txtUnSupportType
		}
		s.reflectObj.Type = message.Text
		return true, txtSetId
	case s.reflectObj.OuterID == "":
		r := []rune(message.Text)
		if len(r) > 15 {
			return true, txtIncorrectId
		}
		s.reflectObj.OuterID = message.Text
		return true, txtSetMinRating
	case s.reflectObj.Settings.MinRating == 0:
		minRating, err := strconv.Atoi(message.Text)
		if err != nil {
			return true, txtErrMinRating
		}
		if minRating < 1 || minRating > 10 {
			return true, txtErrMinRating2
		}
		s.reflectObj.Settings.MinRating = minRating
		return true, txtSetNotLaterDays
	case s.reflectObj.Settings.NotLaterDays == 0:
		notLaterDays, err := strconv.Atoi(message.Text)
		if err != nil {
			return true, txtErrNotLaterDays
		}
		if notLaterDays == 0 {
			return true, txtErrNotLaterDays2
		}
		s.reflectObj.Settings.NotLaterDays = notLaterDays
		return true, txtSetLanguage
	case s.reflectObj.Settings.Language == "":
		if message.Text != "ru" {
			return true, txtErrLanguage
		}
		s.reflectObj.Settings.Language = message.Text
		return true, txtSetPeriod
	case s.Period == nil:
		period, err := strconv.Atoi(message.Text)
		if err != nil {
			return true, txtErrPeriod
		}
		if period < 10 {
			return true, txtErrPeriod2
		}
		s.Period = new(int)
		*s.Period = period
		return true, txtSetPostFormatChoice
	case s.PostToChannels == nil:
		txt := txtSetAllowedFileEmpty
		if message.Text == "0" {
			s.PostToChannels = new(bool)
			*s.PostToChannels = true
			return true, txt
		} else if message.Text == "1" {
			s.PostToChannels = new(bool)
			*s.PostToChannels = false
			return true, txt
		} else {
			return true, txtErrPostToChannels
		}
	case !s.reflectObj.Settings.EmptyFileAllowed:
		if message.Text == "да" {
			s.reflectObj.Settings.EmptyFileAllowed = true
			return false, ""
		} else if message.Text == "нет" {
			s.reflectObj.Settings.EmptyFileAllowed = false
			return false, ""
		} else {
			return true, txtErrAllowedFileEmpty
		}
	default:
		return false, ""
	}
}

type StateActivateCode struct {
	reflectObj *backend.UserEmailActivateRequest
	Access     string
}

func NewStateActivationCode(access string) *StateActivateCode {
	return &StateActivateCode{
		reflectObj: &backend.UserEmailActivateRequest{},
		Access:     access,
	}
}

func (s *StateActivateCode) GetReflectionObject() interface{} {
	return s.reflectObj
}

func (s *StateActivateCode) PerformAdd(message *tgbotapi.Message) (bool, string) {
	switch {
	case s.reflectObj.Code == "":
		s.reflectObj.Code = message.Text
		return false, ""
	default:
		return false, ""
	}
}

type StateAPIChannelSettingsChange struct {
	reflectObj *backend.ChannelSettingsChangeRequest
	Access     string
}

func NewStateAPIChannelSettingsChange(access string) *StateAPIChannelSettingsChange {
	return &StateAPIChannelSettingsChange{
		reflectObj: &backend.ChannelSettingsChangeRequest{
			Data: make(map[string]interface{}),
		},
		Access: access,
	}
}

func (s *StateAPIChannelSettingsChange) GetReflectionObject() interface{} {
	return s.reflectObj
}

func (s *StateAPIChannelSettingsChange) PerformAdd(message *tgbotapi.Message) (bool, string) {
	for x := range s.reflectObj.Data {
		switch x {
		case "min_rating":
			minRating, err := strconv.Atoi(message.Text)
			if err != nil {
				return true, txtErrMinRating
			}
			if minRating < 1 || minRating > 10 {
				return true, txtErrMinRating2
			}
			s.reflectObj.Data[x] = minRating
			break
		case "not_later_days":
			notLaterDays, err := strconv.Atoi(message.Text)
			if err != nil {
				return true, txtErrNotLaterDays
			}
			if notLaterDays == 0 {
				return true, txtErrNotLaterDays2
			}
			s.reflectObj.Data[x] = notLaterDays
			break
		case "empty_file_allowed":
			if message.Text == "да" {
				s.reflectObj.Data[x] = true
				break
			} else if message.Text == "нет" {
				s.reflectObj.Data[x] = false
				break
			} else {
				return true, txtErrAllowedFileEmpty
			}
		case "language":
			if message.Text != "ru" {
				return true, txtErrLanguage
			}
			s.reflectObj.Data[x] = message.Text
			break
		}

		break
	}

	return false, ""
}

type StateClientChannelSettingsChange struct {
	reflectObj map[string]interface{}
	ID         string
	Access     string
}

func NewStateClientAPIChannelSettingsChange(access string) *StateClientChannelSettingsChange {
	return &StateClientChannelSettingsChange{
		reflectObj: make(map[string]interface{}),
		Access:     access,
	}
}

func (s *StateClientChannelSettingsChange) GetReflectionObject() interface{} {
	return s.reflectObj
}

func (s *StateClientChannelSettingsChange) PerformAdd(message *tgbotapi.Message) (bool, string) {
	for x := range s.reflectObj {
		switch x {
		case "period":
			period, err := strconv.Atoi(message.Text)
			if err != nil {
				return true, txtErrPeriod
			}
			if period < 10 {
				return true, txtErrPeriod2
			}
			s.reflectObj[x] = period
			break
		case "post_to_channels":
			if message.Text == "да" {
				s.reflectObj[x] = true
				break
			} else if message.Text == "нет" {
				s.reflectObj[x] = false
				break
			} else {
				return true, txtErrAllowedFileEmpty
			}
		}

		break
	}

	return false, ""
}
