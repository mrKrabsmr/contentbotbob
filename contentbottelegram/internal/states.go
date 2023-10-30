package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
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
