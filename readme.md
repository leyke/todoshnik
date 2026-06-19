
# Тудушник
Приложение сохраняющее задачи в память и отмечающее их статус. Собираю как практику после изучения Golang

## О себе
[Структура](https://github.com/leyke/todoshnik/blob/main/docs/structure.md)

Стэк:
- Go
- REST API
- OpenAPI
- Telegram Bot API
- GORM
- SQL migrations
- Docker / Docker Compose
- gRPC (in progress)
- JSON/File storage
- Context / Graceful shutdown
- Middleware

## Доступные контейнеры

### Консольный
Запуск `go run cmd/cli/main.go`  
[Справочник команд](https://github.com/leyke/todoshnik/blob/main/docs/cli/readme.md)

### АПИ
Запуск `go run cmd/api/main.go`  
[OpenApi спецификция](https://github.com/leyke/todoshnik/blob/main/docs/openapi/openapi.yaml)

### TG Bot
Запуск `go run cmd/tg_bot/main.go`  
[@todoshnik_bot](https://t.me/todoshnik_bot)

### Запуск
#### Установка зависимостей
`make deps`
#### Проверка кода
`make lint`
#### Миграции
##### Применить миграции:
`make migrate-up`
#### Docker
##### Запуск всех сервисов:
`docker compose up --build`