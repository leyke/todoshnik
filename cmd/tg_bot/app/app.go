package app

import (
	"fmt"
	"todoshnik/cmd/tg_bot/app/bot/command"
	"todoshnik/cmd/tg_bot/app/bot/token"
	"todoshnik/internal/config"

	commonapp "todoshnik/internal/app"
	commandstorage "todoshnik/internal/infrastructure/redis/tg_bot/command"
	tokenstorage "todoshnik/internal/infrastructure/redis/tg_bot/token"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	*commonapp.App

	BotApi   *tgbotapi.BotAPI
	Storages *Storages
}

type Storages struct {
	CommandStorage command.Storage
	TokenStorage   token.Storage
}

func InitApp(cfg config.Config) (*App, error) {
	ca, err := commonapp.InitApp(cfg)
	if err != nil {
		return nil, fmt.Errorf("common DI init failed: %w", err)
	}

	botapi, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("create bot API: %w", err)
	}
	botapi.Debug = cfg.Telegram.Debug

	return &App{
		App: ca,

		BotApi: botapi,

		Storages: &Storages{
			CommandStorage: commandstorage.NewStorage(ca.Cache),
			TokenStorage:   tokenstorage.NewStorage(ca.Cache),
		},
	}, nil
}
