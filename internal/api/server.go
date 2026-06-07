package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"todoshnik/internal/app"
	authapi "todoshnik/internal/auth/api"
)

// приложу сюда, хотя к пакету не относится
// посмотри, если не видел пример стандартной структуры проекта https://github.com/golang-standards/project-layout
// может быть наведет на какие-то мысли. В целом мы у себя тоже не придерживаемся строго этого стандарта, но мы следуем
// гексогональной артхитектуре (по крайней мере насколько мы этим владеем и насколько получается)
// В целом есть претензии к структуре.
// Вроде бы есть прекрасные мысли: папки app и infrastructure, но тут же есть cli и client, которые либо сервисы либо
// инфраструктурные (еще не смотрел). Есть какие-то папки намекающие на доменную логику: task, user. Есть богомерзкая
// папка constants, я видел проекты где подобным образом делали, или например метрики/интерфейсы для моков складывали
// в отдельный пакет– это каждый раз плохо заканчивалось и с этим трудно работать.
// - константы ДОЛЖНЫ быть определены в том пакете где они используются, тогда ты а) при удалении пакета
//   или когда реализацию перепишут не забудешь удалить константы б) у тебя не будут 2 разных пакета ссылаться на одну
//   константу.
// - сейчас я вижу что все константы относятся к пакету ui-helper, почему не перенести туда? Кажется что emoji это не
//   то, что можно где-то переиспользовать еще.
// - эмоджи это что-то что относится к слою представления, с этим надо быть особенно осторожным чтобы это нику не
//   протекло, поэтому лучше перенести в пакет ui-helper и сделать неэкспортуемыми (с маленькой буквы)

import (
	taskapi "todoshnik/internal/task/api"
)

type APIHandler struct {
	taskHandler *taskapi.Handler
	authHandler *authapi.Handler
	logger      *log.Logger
	server      *http.Server
}

func NewAPIHandler(container *app.App) *APIHandler {
	return &APIHandler{
		// прекрасно что ты прокидываешь в NewAPIHandler taskHandler, authHandler это прекрасный пример использования
		// DI. Но ты прокидываешь целый application, это слишком жирно и убивает всю идею. Надо прокидывать минимальный
		// уровень зависимостей, например явно 2 хендлера (если они не сгруппированы в какую-то отдельную коллекцию).
		// Если хочется можно на уровне пакета определить api.Handlers структуру и инициализировать ее при вызове
		// _, _ = NewAPIHandler(ctx, logger, handlers)
		// сейчас никто не мешает мне тут начать слать sql запросы, использовать redis и т.п. могу даже шатдаун
		// приложения сделать. Это тяжело и опасно рефакторить
		taskHandler: taskapi.NewHandler(container.TaskService),
		authHandler: authapi.NewHandler(container.UserService, container.TokenService),
		logger:      container.Logger,
	}
}

func (h *APIHandler) Run() error {
	h.server = &http.Server{
		// протекает env и где хост
		Addr:    ":" + os.Getenv("API_PORT"),
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
	fmt.Fprintf(w, "pong")
}
