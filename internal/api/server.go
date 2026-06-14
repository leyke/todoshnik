package api

// приложу сюда, хотя к пакету не относится
// посмотри, если не видел пример стандартной структуры проекта https://github.com/golang-standards/project-layout
// может быть наведет на какие-то мысли. В целом мы у себя тоже не придерживаемся строго этого стандарта, но мы следуем
// гексогональной артхитектуре (по крайней мере насколько мы этим владеем и насколько получается)
// В целом есть претензии к структуре.
// Вроде бы есть прекрасные мысли: папки app и infrastructure, но тут же есть cli и client, которые либо сервисы либо
// инфраструктурные (еще не смотрел). Есть какие-то папки намекающие на доменную логику: task, user.

import (
	"context"
	"log"
	"net/http"

	"todoshnik/internal/api/response"
	"todoshnik/internal/app"
	"todoshnik/internal/config"

	authapi "todoshnik/internal/auth/api"
	taskapi "todoshnik/internal/task/api"
)

type APIHandler struct {
	taskHandler *taskapi.Handler
	authHandler *authapi.Handler
	logger      *log.Logger
	server      *http.Server

	config Config
}

func NewAPIHandler(services *app.Services, logger *log.Logger, config Config) *APIHandler {
	return &APIHandler{
		taskHandler: taskapi.NewHandler(services.TaskService),
		authHandler: authapi.NewHandler(services.UserService, services.TokenService, logger),
		logger:      logger,
		config:      config,
	}
}

func (h *APIHandler) Run(cfg config.AppConfig) error {
	h.server = &http.Server{
		Addr:    ":" + cfg.Port, // в докере с хостом не работает
		Handler: h.Router(),
	}

	h.logger.Println("API hello")

	err := h.server.ListenAndServe()
	// 0. Вопрос: ListenAndServe блокирующий вызов или нет? Как можно сделать по-другому?
	// 1. ошибку лучше заворачивать в конкретные ошибки, иначе тяжело понять какой компонент за это отвечает,
	// обязательно это изучи: fmt.Errorf(%w) и error.Join кажется, про джоин отдельная история, им надо пользоваться
	// в конкретных ситуациях для которых он придуман
	// 2. почему ты отдельно обрабатываешь http.ErrServerClosed и подавляешь эту ошибку? Разве это не проблема когда
	// ты пытаешься запустить закрытый/непроинициализированных сервер?
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (h *APIHandler) Shutdown(ctx context.Context) error {
	h.logger.Println("отключение API...")
	return h.server.Shutdown(ctx)
}

func (h *APIHandler) pingHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, "pong")
}
