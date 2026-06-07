package bot

import (
	"log"
	"os"
	"sync"
	"time"

	"todoshnik/internal/app"
	authbot "todoshnik/internal/auth/bot"
	"todoshnik/internal/bot/session"
	"todoshnik/internal/client"
	task "todoshnik/internal/task/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

const botTokenTtl time.Duration = 4 * time.Hour

type Handler struct {
	taskHandler *task.Handler
	authHandler *authbot.Handler

	sessionStorage *session.Storage

	bot *tgbotapi.BotAPI
	// очень смущает что тут на сервисном слое редис клиент; Хендлер не должен работать напрямую с инфраструктурныи
	// зависимостями
	cache *redis.Client
	wg    sync.WaitGroup

	logger *log.Logger
}

func NewHandler(container *app.App, bot *tgbotapi.BotAPI) *Handler {
	apiClient := client.NewApiClient(os.Getenv("API_URL"))
	return &Handler{
		taskHandler:    task.NewHandler(apiClient, container.Logger),
		authHandler:    authbot.NewHandler(apiClient),
		sessionStorage: session.NewStorage(container.Cache),
		cache:          container.Cache,

		bot:    bot,
		logger: container.Logger,
	}
}

func (h *Handler) Run() error {
	log.Printf("Authorized on account %s", h.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.wg.Add(1)

		go func(update tgbotapi.Update) {
			defer h.wg.Done()

			h.dispatchUpdate(update)
		}(update)
	}

	// 1. не хватает wait group возможна или утечка или частичный резултат
	//    вернее группа есть, но без wg.Wait() это бессмыслено, разберись как это работает, поэксперементируй
	// 2. если updates может быть большим– стоит ограничить параллелизм семафором
	// 3. Изучи пакет errorGroup но постарайся сначала сам написать его с нуля

	return nil
}

func (h *Handler) Shutdown() {
	h.logger.Println("Остановка ТГ...")

	h.bot.StopReceivingUpdates()

	h.wg.Wait()

	h.logger.Println("Чтение тг остановлено")
}
