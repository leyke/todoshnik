package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const defaultTokenTTL time.Duration = 7 * 24 * time.Hour
const defaultApiRequestTimeout time.Duration = 5 * time.Second
const defaultSemaphoreSize int = 200

const defaultDbMaxOpenConnections = 25
const defaultDbMaxIdleConnections = 10

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
	Host        string
	Port        string
	TokenSecret string
	TokenTtl    time.Duration
	TmpDir      string
}

type PostgresConfig struct {
	Host                 string
	Port                 string
	User                 string
	Password             string
	DBName               string
	SSLMode              string
	DbMaxOpenConnections int
	DbMaxIdleConnections int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type TelegramConfig struct {
	SemaphoreSize     int
	Token             string
	ServiceToken      string
	BotApiUrl         string
	ApiRequestTimeOut time.Duration
	Debug             bool
}

func Load() (Config, error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	tokenTtl, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		tokenTtl = defaultTokenTTL
	}

	apiRequestTimeOut, err := time.ParseDuration(os.Getenv("API_REQUEST_TIMEOUT"))
	if err != nil {
		tokenTtl = defaultApiRequestTimeout
	}

	semaphoreSize, err := strconv.Atoi(os.Getenv("BOT_SEMAPHORE_SIZE"))
	if err != nil {
		semaphoreSize = defaultSemaphoreSize
	}

	dbMaxOpenConnections, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		semaphoreSize = defaultDbMaxOpenConnections
	}

	dbMaxIdleConnections, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		semaphoreSize = defaultDbMaxIdleConnections
	}

	return Config{
		App: AppConfig{
			Host:        os.Getenv("API_HOST"),
			Port:        os.Getenv("API_PORT"),
			TokenSecret: os.Getenv("SECRET"),
			TokenTtl:    tokenTtl,
			TmpDir:      os.Getenv("TMP_DIR"),
		},
		Postgres: PostgresConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),

			DbMaxOpenConnections: dbMaxOpenConnections,
			DbMaxIdleConnections: dbMaxIdleConnections,
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		Telegram: TelegramConfig{
			SemaphoreSize:     semaphoreSize,
			Token:             os.Getenv("TELEGRAM_TOKEN"),
			ServiceToken:      os.Getenv("BOT_SERVICE_TOKEN"),
			BotApiUrl:         os.Getenv("BOT_API_URL"),
			ApiRequestTimeOut: apiRequestTimeOut,
			Debug:             os.Getenv("TELEGRAM_DEBUG") == "1",
		},
	}, nil
}
