// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package project

import (
	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/managers"
	"github.com/roman-kart/go-initial-project/project/tools"
	"github.com/roman-kart/go-initial-project/project/utils"
)

// Injectors from wire.go:

func InitializeApplication(configFolder string) (*Application, func(), error) {
	configConfig, err := config.NewConfig(configFolder)
	if err != nil {
		return nil, nil, err
	}
	logger, cleanup, err := utils.NewLogger(configConfig)
	if err != nil {
		return nil, nil, err
	}
	errorWrapperCreator := tools.NewErrorWrapperCreator()
	clickHouse, cleanup2, err := utils.NewClickHouse(configConfig, logger, errorWrapperCreator)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	postgresql, cleanup3, err := utils.NewPostgresql(configConfig, logger, errorWrapperCreator)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	rabbitMQ := utils.NewRabbitMQ(configConfig, logger, errorWrapperCreator)
	s3 := utils.NewS3(configConfig, logger, postgresql, errorWrapperCreator)
	telegramBot := utils.NewTelegram(configConfig, logger, rabbitMQ, errorWrapperCreator)
	statManager, err := managers.NewStatManager(logger, clickHouse, configConfig, errorWrapperCreator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	userAccountManager, err := managers.NewUserAccountManager(logger, postgresql, configConfig, errorWrapperCreator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	telegramBotManager, cleanup4, err := managers.NewTelegramBotManager(configConfig, logger, statManager, userAccountManager, telegramBot, errorWrapperCreator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	application := NewApplication(configConfig, clickHouse, logger, postgresql, rabbitMQ, s3, telegramBot, statManager, telegramBotManager, userAccountManager)
	return application, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
