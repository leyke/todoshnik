# Тудушник

`todoshnik` — Go-приложение для управления задачами.

Проект создан как практика изучения Go и построен в стиле Hexagonal / Clean Architecture.

Основные компоненты:

- REST API;
- Telegram Bot;
- CLI.

## Документация

- [Архитектура проекта](docs/structure.md)
- [CLI команды](docs/cli/readme.md)
- [OpenAPI спецификация](docs/openapi/openapi.yaml)
- [Команды запуска и Makefile](docs/makefile.md)

## Стек

- Go 1.26.4
- PostgreSQL 16
- Redis
- Chi
- Goose migrations
- Docker / Docker Compose
- OpenAPI
- Telegram Bot API
- Squirrel
- GORM (для примера) 

## Компоненты

### API

REST API для:

- регистрации и авторизации;
- управления задачами;
- Telegram-аутентификации.

Спецификация: [OpenAPI](docs/openapi/openapi.yaml)

---

### Telegram Bot

Работа с задачами через Telegram.

[@todoshnik_bot](https://t.me/todoshnik_bot)

---

### CLI

Консольный интерфейс для управления задачами.

Документация: [CLI](docs/cli/readme.md)

---

### Миграции

Отдельный executable для применения миграций PostgreSQL через Goose.

---

## Запуск

Все команды запуска, сборки и проверки: [docs/makefile.md](docs/makefile.md)

QuickStart:

```bash
docker compose up --build