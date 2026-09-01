# Архитектура проекта todoshnik

## Общее описание

`todoshnik` — Go-приложение для управления задачами с несколькими точками взаимодействия:

- REST API;
- Telegram Bot;
- CLI интерфейс.

Проект построен с использованием подхода, близкого к Hexagonal Architecture / Clean Architecture:

- бизнес-логика находится в доменных слоях;
- инфраструктура вынесена отдельно;
- зависимости передаются через Dependency Injection;
- домены работают через интерфейсы;
- конкретные реализации подключаются на уровне приложения.

Основные технологии:

- Go 1.26.4;
- PostgreSQL 16;
- Redis;
- REST API на chi;
- Telegram Bot API;
- Docker;
- Goose migrations.

---

# Общая структура каталогов

```text
.
├── bin
│   ├── golangci-lint
│   └── goose
├── cmd
│   ├── api
│   │   ├── app
│   │   │   ├── auth
│   │   │   ├── errors
│   │   │   ├── middleware
│   │   │   ├── request
│   │   │   ├── response
│   │   │   ├── task
│   │   │   ├── config.go
│   │   │   ├── router.go
│   │   │   └── server.go
│   │   └── main.go
│   ├── cli
│   │   ├── app
│   │   │   ├── add.go
│   │   │   ├── clear_tokens.go
│   │   │   ├── delete.go
│   │   │   ├── done.go
│   │   │   ├── handler.go
│   │   │   ├── helper.go
│   │   │   ├── list.go
│   │   │   └── printer.go
│   │   └── main.go
│   ├── migrate
│   │   └── main.go
│   └── tg_bot
│       ├── app
│       │   ├── auth
│       │   ├── bot
│       │   ├── task
│       │   └── app.go
│       └── main.go
├── docker
│   ├── api.Dockerfile
│   ├── api.debug.Dockerfile
│   ├── bot.Dockerfile
│   └── migrate.Dockerfile
├── docs
│   ├── cli
│   ├── openapi
│   └── structure.md
├── internal
│   ├── app
│   │   ├── app.go
│   │   └── logger.go
│   ├── config
│   │   └── config.go
│   ├── domains
│   │   ├── task
│   │   ├── token
│   │   └── user
│   └── infrastructure
│       ├── api_client
│       ├── context_manager
│       ├── db
│       ├── filestorage
│       ├── identity
│       ├── redis
│       ├── security
│       ├── utils
│       └── validation
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Описание каталогов

- `cmd` — точки входа приложения. Содержит отдельные executable:
  - `api` — HTTP REST API сервер, маршруты, middleware, обработчики запросов и ответов.
  - `tg_bot` — Telegram-бот, обработка команд, взаимодействие с API и Redis-хранилищами.
  - `cli` — консольный интерфейс для управления задачами и токенами.
  - `migrate` — утилита запуска миграций базы данных через Goose.

- `docker` — Dockerfile для сборки и запуска компонентов приложения:
  - `api.Dockerfile` — production-сборка API сервера.
  - `api.debug.Dockerfile` — debug-сборка API с поддержкой Delve.
  - `bot.Dockerfile` — сборка Telegram-бота.
  - `migrate.Dockerfile` — контейнер для применения миграций.

- `docker-compose.yml` — конфигурация запуска полного окружения:
  - API;
  - Telegram-бот;
  - PostgreSQL;
  - Redis;
  - миграции базы данных.

- `docs` — документация проекта:
  - `openapi` — спецификация REST API в формате OpenAPI.
  - `cli` — документация по CLI интерфейсу.
  - `structure.md` — описание архитектуры и структуры проекта.

- `internal` — внутренняя реализация приложения. Код из данного пакета не предназначен для использования внешними модулями.

  - `app` — сборка приложения и Dependency Injection:
    - создание общих сервисов;
    - настройка logger;
    - управление жизненным циклом приложения.

  - `config` — загрузка и описание конфигурации приложения:
    - API настройки;
    - PostgreSQL;
    - Redis;
    - Telegram Bot.

  - `domains` — бизнес-логика приложения, разделенная по предметным областям:
    
    - `task` — доменная логика задач:
      - модели задач;
      - сервисы;
      - фильтры;
      - интерфейс репозитория.

    - `user` — доменная логика пользователей:
      - модель пользователя;
      - сервисы;
      - работа с паролями;
      - интерфейс репозитория.

    - `token` — доменная логика токенов:
      - генерация токенов;
      - проверка сроков действия;
      - управление токенами;
      - интерфейсы репозитория и генератора.

  - `infrastructure` — реализации внешних зависимостей и технические адаптеры:

    - `api_client` — клиент для взаимодействия с внешним или внутренним API.

    - `context_manager` — управление контекстом приложения:
      - хранение данных авторизованного пользователя;
      - передача auth context между слоями.

    - `db` — работа с базой данных:
      - подключение к PostgreSQL;
      - GORM конфигурация;
      - миграции;
      - реализации репозиториев;
      - управление транзакциями.

      Подкаталоги:
      - `repository/task` — хранение и получение задач.
      - `repository/user` — работа с пользователями.
      - `repository/token` — работа с токенами.
      - `transaction` — управление транзакциями.

    - `filestorage` — файловое хранилище и интерфейсы работы с ним.

    - `identity` — управление правами доступа и областью видимости пользователя.

    - `redis` — работа с Redis:
      - подключение;
      - Telegram Bot storage:
        - command storage;
        - token storage.

    - `security` — компоненты безопасности:
      - хеширование паролей;
      - работа с токенами.

    - `utils` — вспомогательные инструменты:
      - работа с временем;
      - тестовые утилиты.

    - `validation` — общие правила валидации данных и обработка ошибок валидации.

- `bin` — локальные бинарные инструменты проекта:
  - golangci-lint;
  - goose.

- `go.mod` — описание Go-модуля и зависимостей проекта.

- `go.sum` — контрольные суммы зависимостей Go.