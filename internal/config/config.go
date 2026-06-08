package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Telegram TelegramConfig
}

type AppConfig struct {
	Host         string
	Port         string
	TokenSalt    string
	TokenTtlDays string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type TelegramConfig struct {
	Token string
	Debug bool
}

func Load() (*Config, error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	return &Config{
		App: AppConfig{
			Host:         os.Getenv("APP_HOST"),
			Port:         os.Getenv("APP_PORT"),
			TokenSalt:    os.Getenv("SALT"),
			TokenTtlDays: os.Getenv("TOKEN_TTL_DAYS"),
		},
		Postgres: PostgresConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		Telegram: TelegramConfig{
			Token: os.Getenv("TELEGRAM_TOKEN"),
			Debug: os.Getenv("TELEGRAM_DEBUG") == "1",
		},
	}, nil
}
