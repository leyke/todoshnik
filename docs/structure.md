# Структура проекта

Проект `todoshnik` организован по классической архитектуре Go-приложения с отдельными точками входа, бизнес-логикой и вспомогательными компонентами.

```text
.
├── cmd
│   ├── api
│   │   └── main.go
│   ├── cli
│   │   └── main.go
│   └── tg_bot
│       └── main.go
├── docker
│   ├── api.Dockerfile
│   └── bot.Dockerfile
├── docker-compose.yml
├── docs
│   ├── cli
│   │   └── readme.md
│   ├── openapi
│   │   └── openapi.yaml
|   └── structure.md
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── middleware
│   │   │   ├── auth.go
│   │   │   └── logging.go
│   │   ├── request
│   │   │   └── json.go
│   │   ├── response
│   │   │   ├── error.go
│   │   │   └── json.go
│   │   ├── router.go
│   │   └── server.go
│   ├── app
│   │   ├── app.go
│   │   └── logger.go
│   ├── auth
│   │   ├── api
│   │   │   ├── dto.go
│   │   │   ├── get_authorized_user.go
│   │   │   ├── handler.go
│   │   │   ├── sign_in.go
│   │   │   ├── sign_up.go
│   │   │   ├── tg_auto_reg.go
│   │   │   └── tg_login.go
│   │   ├── bot
│   │   │   ├── dto.go
│   │   │   ├── get_token.go
│   │   │   ├── handler.go
│   │   │   └── sign_in_user.go
│   │   ├── config.go
│   │   ├── context
│   │   │   ├── context.go
│   │   │   ├── token.go
│   │   │   └── user.go
│   │   ├── device.go
│   │   ├── helper.go
│   │   └── token
│   │       ├── filestorage.go
│   │       ├── helper.go
│   │       ├── repository.go
│   │       ├── repository_db.go
│   │       ├── repository_file.go
│   │       ├── service.go
│   │       └── token.go
│   ├── bot
│   │   ├── auth.go
│   │   ├── callback.go
│   │   ├── command.go
│   │   ├── handler.go
│   │   ├── message.go
│   │   ├── response
│   │   │   ├── error.go
│   │   │   ├── inline_keyboard.go
│   │   │   └── task.go
│   │   ├── sender.go
│   │   ├── session
│   │   │   ├── session.go
│   │   │   ├── state_helper.go
│   │   │   └── storage.go
│   │   ├── tg
│   │   │   ├── callback.go
│   │   │   ├── command.go
│   │   │   └── state.go
│   │   └── update_dispatcher.go
│   ├── cli
│   │   ├── add.go
│   │   ├── clear_tokens.go
│   │   ├── delete.go
│   │   ├── done.go
│   │   ├── handler.go
│   │   ├── helper.go
│   │   ├── list.go
│   │   └── printer.go
│   ├── client
│   │   └── client_api.go
│   ├── constants
│   │   └── emoji.go
│   ├── errors
│   │   └── errors.go
│   ├── identity
│   │   └── scope.go
│   ├── infrastructure
│   │   ├── db
│   │   │   ├── gorm.go
│   │   │   └── migrations
│   │   └── redis
│   │       └── redis.go
│   ├── storage
│   │   ├── filestorage.go
│   │   └── storage.go
│   ├── task
│   │   ├── api
│   │   │   ├── create.go
│   │   │   ├── delete.go
│   │   │   ├── done.go
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── helpers.go
│   │   │   ├── list.go
│   │   │   ├── update.go
│   │   │   └── view.go
│   │   ├── bot
│   │   │   ├── add.go
│   │   │   ├── delete.go
│   │   │   ├── done.go
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── list.go
│   │   │   └── ui_helper.go
│   │   ├── filestorage.go
│   │   ├── filter.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── repository_db.go
│   │   ├── repository_file.go
│   │   └── service.go
│   ├── user
│   │   ├── filestorage.go
│   │   ├── repository.go
│   │   ├── repository_db.go
│   │   ├── repository_file.go
│   │   ├── service.go
│   │   └── user.go
│   └── validation
│       └── validation.go
├── readme.md
└── tmp
    ├── api.log
    ├── cli.log
    ├── tasks.json
    ├── tg.log
    ├── tokens.json
    └── users.json
```

Описание каталогов:

- `cmd` — точки входа для запуска API-сервера, CLI-клиента и Telegram-бота.
- `docker` — Dockerfile для контейнеризации API и бота.
- `docker-compose.yml` — конфигурация для запуска сервисов вместе.
- `docs` — документация по CLI и OpenAPI спецификация.
- `internal` — внутренняя реализация приложения, недоступная для внешних пакетов.
  - `api` — HTTP-сервер, маршруты, middleware и JSON-обработка.
  - `app` — инициализация приложения и логирование.
  - `auth` — авторизация, работа с токенами, контекстом и Telegram-аутентификацией.
  - `bot` — логика Telegram-бота, обработка команд, сессий и ответов.
  - `cli` — локальный интерфейс командной строки для управления задачами.
  - `client` — API-клиент для внутреннего или внешнего использования.
  - `constants` — константы, в том числе эмодзи.
  - `errors` — общее определение ошибок.
  - `identity` — управление правами доступа и областью видимости.
  - `infrastructure` — адаптеры для работы с БД и Redis.
  - `storage` — файловое хранилище и общий интерфейс хранения.
  - `task` — бизнес-логика работы с задачами для API, бота и хранилища.
  - `user` — хранение и сервисы пользователей.
  - `validation` — общие правила валидации данных.
- `tmp` — временные файлы и логи, используемые во время запуска.
- `readme.md`, `go.mod`, `go.sum` — документация проекта и зависимости Go.
