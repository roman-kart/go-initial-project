package project

import (
	"go.uber.org/zap"

	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/managers"
	"github.com/roman-kart/go-initial-project/project/utils"
)

// Application contains all components of go-initial-project components.
type Application struct {
	Config *config.Config

	ClickHouse  *utils.ClickHouse
	Logger      *utils.Logger
	logger      *zap.Logger
	Postgres    *utils.Postgresql
	RabbitMQ    *utils.RabbitMQ
	S3          *utils.S3
	TelegramBot *utils.TelegramBot

	StatManager        *managers.StatManager
	TelegramBotManager *managers.TelegramBotManager
	UserAccountManager *managers.UserAccountManager
}

// NewApplication creates a new instance of Application.
// Using for configuring with wire.
func NewApplication(
	cfg *config.Config,

	clickHouse *utils.ClickHouse,
	logger *utils.Logger,
	postgres *utils.Postgresql,
	rabbitmq *utils.RabbitMQ,
	s3 *utils.S3,
	telegramBot *utils.TelegramBot,

	statManager *managers.StatManager,
	telegramBotManager *managers.TelegramBotManager,
	userAccountManager *managers.UserAccountManager,
) *Application {
	return &Application{
		Config: cfg,

		ClickHouse:  clickHouse,
		Logger:      logger,
		logger:      logger.Logger.Named("Application"),
		Postgres:    postgres,
		RabbitMQ:    rabbitmq,
		S3:          s3,
		TelegramBot: telegramBot,

		StatManager:        statManager,
		TelegramBotManager: telegramBotManager,
		UserAccountManager: userAccountManager,
	}
}
