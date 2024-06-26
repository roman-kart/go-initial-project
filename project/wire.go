//go:build wireinject
// +build wireinject

package project

import (
	"github.com/google/wire"
	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/managers"
	"github.com/roman-kart/go-initial-project/project/tools"
	"github.com/roman-kart/go-initial-project/project/utils"
)

func InitializeApplication(configFolder string) (*Application, func(), error) {
	wire.Build(
		config.NewConfig,

		tools.NewErrorWrapperCreator,
		utils.NewClickHouse,
		utils.NewLogger,
		utils.NewPostgresql,
		utils.NewRabbitMQ,
		utils.NewS3,
		utils.NewTelegram,

		managers.NewStatManager,
		managers.NewTelegramBotManager,
		managers.NewUserAccountManager,

		NewApplication,
	)
	return &Application{}, func() {}, nil
}
