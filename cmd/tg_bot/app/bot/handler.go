package bot

import (
	"log"
	"sync"
	"time"

	"todoshnik/cmd/tg_bot/app"
	"todoshnik/cmd/tg_bot/app/bot/command"
	"todoshnik/cmd/tg_bot/app/bot/token"

	authbot "todoshnik/cmd/tg_bot/app/auth"
	taskbot "todoshnik/cmd/tg_bot/app/task"
	client "todoshnik/internal/infrastructure/api_client"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	semaphore   chan struct{}
	taskHandler *taskbot.Handler
	authHandler *authbot.Handler

	commandStorage command.Storage
	tokenStorage   token.Storage

	bot *tgbotapi.BotAPI

	wg     sync.WaitGroup
	logger *log.Logger
}

type Config struct {
	ApiURL            string
	ApiRequestTimeout time.Duration
	BotServiceToken   string
	SemaphoreSize     int
}

func NewHandler(container *app.App, cfg Config) *Handler {
	apiClient := client.NewApiClient(cfg.ApiURL, cfg.BotServiceToken, cfg.ApiRequestTimeout)
	container.Logger.Printf("new BOT handlerApiRequestTimeout: %v", cfg.ApiRequestTimeout)
	return &Handler{
		semaphore:   make(chan struct{}, cfg.SemaphoreSize),
		taskHandler: taskbot.NewHandler(apiClient, container.Logger),
		authHandler: authbot.NewHandler(
			apiClient,
			authbot.Config{
				RequestTimeout: cfg.ApiRequestTimeout,
			},
			container.Logger,
		),
		commandStorage: container.Storages.CommandStorage,
		tokenStorage:   container.Storages.TokenStorage,

		bot:    container.BotApi,
		logger: container.Logger,
	}
}

func (h *Handler) Run() error {
	h.logger.Printf("Authorized on account %s", h.bot.Self.UserName)
	h.semaphore <- struct{}{}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.wg.Add(1)

		go func(update tgbotapi.Update) {
			defer h.wg.Done()
			defer func() { <-h.semaphore }()

			h.dispatchUpdate(update)
		}(update)
	}

	h.wg.Wait()

	return nil
}

func (h *Handler) Shutdown() {
	h.logger.Println("Остановка ТГ...")

	h.bot.StopReceivingUpdates()

	h.wg.Wait()

	h.logger.Println("Чтение тг остановлено")
}
