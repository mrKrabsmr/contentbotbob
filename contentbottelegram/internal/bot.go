package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/backend"
	"github.com/mrKrabsmr/contentbottelegram/internal/clients/redis"
	"github.com/mrKrabsmr/contentbottelegram/internal/configs"
	"github.com/sirupsen/logrus"
	"time"
)

type Bot struct {
	config     *configs.Config
	logger     *logrus.Logger
	bot        *tgbotapi.BotAPI
	redis      *redis.Client
	api        *backend.APIClient
	storage    map[int64]State
	stopManage map[string]chan struct{}
}

func NewBot(config *configs.Config) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	return &Bot{
		config:     config,
		logger:     logger,
		bot:        botAPI,
		redis:      redis.NewClient(config.Redis),
		api:        backend.NewAPIClient(config.Backend, logger),
		storage:    make(map[int64]State),
		stopManage: make(map[string]chan struct{}),
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

func (b *Bot) manageRun() error {
	channels, err := b.api.ChannelList()
	if err != nil {
		return err
	}

	for _, channel := range channels {
		chatId, err := b.redis.GetPostChatId(channel.OuterID)
		if err != nil {
			continue
		}
		
		period, err := b.redis.GetPostPeriod(channel.OuterID)
		if err != nil {
			continue
		}

		go b.StartManage(channel.OuterID, chatId, period)
	}

	return nil
}

func (b *Bot) Run() error {
	if err := b.configure(); err != nil {
		return err
	}

	time.Sleep(time.Second * 5)
	if err := b.manageRun(); err != nil {
		b.logger.Error(err)
	}

	u := tgbotapi.NewUpdate(0)
	updates := b.bot.GetUpdatesChan(u)

	b.handleUpdates(updates)

	return nil
}
