package managers

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/telebot.v3"

	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/tools"
	"github.com/roman-kart/go-initial-project/project/utils"
)

// TelegramBotManager managing [utils.TelegramBot].
type TelegramBotManager struct {
	Config              *config.Config
	Logger              *utils.Logger
	logger              *zap.Logger
	TelegramBot         *utils.TelegramBot
	telegramBot         *telebot.Bot
	StatManager         *StatManager
	UserAccountManager  *UserAccountManager
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewTelegramBotManager creates new TelegramBotManager.
// Using for configuring with wire.
func NewTelegramBotManager(
	config *config.Config,
	logger *utils.Logger,
	statManager *StatManager,
	userAccountManager *UserAccountManager,
	telegramBot *utils.TelegramBot,
	errorWrapperCreator tools.ErrorWrapperCreator,
) (*TelegramBotManager, func(), error) {
	tbm := &TelegramBotManager{
		Config:              config,
		Logger:              logger,
		logger:              logger.Logger.Named("TelegramBotManager"),
		TelegramBot:         telegramBot,
		StatManager:         statManager,
		UserAccountManager:  userAccountManager,
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("TelegramBotManager"),
	}

	ew := tools.GetErrorWrapper("NewTelegramBotManager")

	err := tbm.createBot()
	if err != nil {
		return nil, nil, ew(err)
	}

	go tbm.telegramBot.Start()

	return tbm, func() { tbm.telegramBot.Stop() }, nil
}

func (t *TelegramBotManager) createBot() error {
	ew := t.ErrorWrapperCreator.GetMethodWrapper("createBot")

	bot, err := t.TelegramBot.CreateBotDefault()
	if err != nil {
		return ew(err)
	}

	t.telegramBot = bot

	return nil
}

// GetBot returns [utils.TelegramBot] instance.
// If bot is not created, it will be created, but will panic if error occurred.
func (t *TelegramBotManager) GetBot() *telebot.Bot {
	if t.telegramBot == nil {
		tools.PanicOnError(t.createBot())
	}

	return t.telegramBot
}

// StartCommandConfig contains configurations for start command.
type StartCommandConfig struct {
	Enabled bool
	Message string
}

// HelpCommandMessages contains configurations for help command.
type HelpCommandMessages struct {
	ShortMessage  string
	DetailMessage string
}

// HelpCommandConfig contains configuration of one command's help message.
type HelpCommandConfig struct {
	Enabled              bool
	MainHelpMessage      string
	CommandsHelpMessages map[string]HelpCommandMessages
}

// CommonBotCommandsConfig contains configurations for common bot commands.
type CommonBotCommandsConfig struct {
	Start StartCommandConfig
	Help  HelpCommandConfig
}

// ErrNoMessage error if no message provided for command.
var ErrNoMessage = errors.New("no message for start command")

// AddCommonCommandsHandlers adds handlers for common bot commands.
//
//nolint:funlen
func (t *TelegramBotManager) AddCommonCommandsHandlers(cfg *CommonBotCommandsConfig) {
	if cfg.Start.Enabled {
		ew := t.ErrorWrapperCreator.GetMethodWrapper("start_handler")

		t.GetBot().Handle("/start", func(c telebot.Context) error {
			if cfg.Start.Message != "" {
				return ew(c.Send(cfg.Start.Message, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown}))
			}

			return ew(ErrNoMessage)
		})
	}

	if cfg.Help.Enabled {
		ew := t.ErrorWrapperCreator.GetMethodWrapper("help_handler")

		t.GetBot().Handle("/help", func(c telebot.Context) error {
			if cfg.Help.MainHelpMessage == "" {
				return ew(ErrNoMessage)
			}

			args := c.Args()
			if len(args) == 0 {
				commandsListMessagePart := ""

				for commandName, commandConfig := range cfg.Help.CommandsHelpMessages {
					commandsListMessagePart += fmt.Sprintf("%s - %s\n", commandName, commandConfig.ShortMessage)
				}

				finalMessage := fmt.Sprintf("%s\n\n*Команды:*\n%s", cfg.Help.MainHelpMessage, commandsListMessagePart)

				return ew(c.Send(finalMessage, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown}))
			}

			commandName := args[0]
			commandName = strings.TrimSpace(commandName)

			// remove leading slash if it exists for make sure only one slash will exists
			commandName = "/" + strings.TrimLeft(commandName, "/")

			commandConfig, ok := cfg.Help.CommandsHelpMessages[commandName]
			if !ok {
				return ew(
					c.Send(
						fmt.Sprintf("Команда `%s` не найдена", commandName),
						&telebot.SendOptions{
							ParseMode: telebot.ModeMarkdown,
						},
					),
				)
			}

			return ew(
				c.Send(
					fmt.Sprintf("%s\n\n%s", commandConfig.ShortMessage, commandConfig.DetailMessage),
					&telebot.SendOptions{
						ParseMode: telebot.ModeMarkdown,
					},
				),
			)
		})
	}
}