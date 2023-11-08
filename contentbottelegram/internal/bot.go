package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/redis"
	"github.com/mrKrabsmr/contentbottelegram/internal/configs"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	config  *configs.Config
	logger  *logrus.Logger
	bot     *tgbotapi.BotAPI
	redis   *redis.Client
	api     *backend.APIClient
	storage map[int64]State
}

func NewBot(config *configs.Config) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	return &Bot{
		config:  config,
		logger:  logger,
		bot:     botAPI,
		redis:   redis.NewClient(config.Redis),
		api:     backend.NewAPIClient(config.Backend, logger),
		storage: make(map[int64]State),
	}, nil
}

func (b *Bot) configure() error {
	// configure logger

	logLevel, err := logrus.ParseLevel(b.config.LogLevel)
	if err != nil {
		return err
	}

	b.logger.SetLevel(logLevel)

	// configure bot

	if b.config.Debug {
		b.bot.Debug = true
	}

	return nil
}

func (b *Bot) Run() error {
	if err := b.configure(); err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	updates := b.bot.GetUpdatesChan(u)

	b.handleUpdates(updates)

	return nil
}
