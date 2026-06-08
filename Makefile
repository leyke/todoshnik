include .env
export

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up