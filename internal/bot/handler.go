package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"todoshnik/internal/app"
	authbot "todoshnik/internal/auth/bot"
	"todoshnik/internal/bot/session"
	"todoshnik/internal/bot/tg"
	"todoshnik/internal/client"
	apperrors "todoshnik/internal/errors"
	task "todoshnik/internal/task/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

const botTokenTtl time.Duration = 4 * time.Hour

type BotHandler struct {
	TaskHandler  *task.Handler
	AuthHandler  *authbot.Handler
	StateStorage *session.StateStorage
	cache        *redis.Client
	api          *client.ApiClient
	bot          *tgbotapi.BotAPI
	logger       *log.Logger
}

func NewBotHandler(container *app.App, bot *tgbotapi.BotAPI) *BotHandler {
	apiClient := client.NewApiClient(os.Getenv("API_URL"))
	return &BotHandler{
		TaskHandler:  task.NewHandler(apiClient, container.Logger),
		AuthHandler:  authbot.NewHandler(apiClient),
		StateStorage: session.NewStateStorage(container.Cache),
		cache:        container.Cache,
		api:          apiClient,
		bot:          bot,
		logger:       container.Logger,
	}
}

func (bh *BotHandler) Run() {
	log.Printf("Authorized on account %s", bh.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bh.bot.GetUpdatesChan(u)

	for update := range updates {
		go bh.handleUpdate(update)
	}
}

func (bh *BotHandler) handleUpdate(update tgbotapi.Update) {
	var msg *tgbotapi.MessageConfig

	ctx := context.Background()

	if update.CallbackQuery != nil {
		ctx = bh.handleAuth(ctx, update.CallbackQuery.From)
		msg = bh.handleCallback(ctx, update)
	} else if update.Message != nil {
		ctx = bh.handleAuth(ctx, update.Message.From)
		if update.Message.IsCommand() {
			msg = bh.handleCommand(ctx, update)
		} else {
			msg = bh.handleMessage(ctx, update)
		}
	} else {
		return
	}

	if _, err := bh.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (bh *BotHandler) handleCallback(ctx context.Context, update tgbotapi.Update) *tgbotapi.MessageConfig {
	query := update.CallbackQuery
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")

	bh.bot.Request(tgbotapi.NewCallback(query.ID, ""))

	var callback *tg.CallbackQuery
	err := json.Unmarshal([]byte(query.Data), &callback)
	if err != nil {
		msg.Text = "Возникла непредвиденная ошибка"
		fmt.Println(query.Data)
		fmt.Println(err.Error())
		bh.logger.Println(err)
		return &msg
	}

	switch callback.Command {
	case tg.СommandTaskDone:
		taskRowText, err := bh.TaskHandler.DoneTask(ctx, callback.Payload["task_id"])
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				msg.Text = err.Error()
			} else if errors.Is(err, apperrors.ErrUnAuth) {
				msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
				return &msg
			} else {
				msg.Text = "Возникла непредвиденная ошибка"
				fmt.Println(err.Error())
				bh.logger.Println(err)
			}
			break
		}
		if taskRowText != "" {
			bh.editTgMessage(query.Message.Chat.ID, query.Message.MessageID, taskRowText)
		}

		msg.Text = "Статус обновлен"
	case tg.CommandTaskDelete:
		err := bh.TaskHandler.DeleteTask(ctx, callback.Payload["task_id"])
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				msg.Text = err.Error()
			} else if errors.Is(err, apperrors.ErrUnAuth) {
				msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
				return &msg
			} else {
				msg.Text = "Возникла непредвиденная ошибка"
				fmt.Println(err.Error())
				bh.logger.Println(err)
			}
			break
		}
		fmt.Println(query)
		msg.Text = "Задача удалена"
		bh.deleteTgMessage(query.Message.Chat.ID, query.Message.MessageID)
	default:
		msg.Text = "Я хз что это такое, если бы я знал что это такое, я бы помог"
	}

	return &msg
}

func (bh *BotHandler) editTgMessage(chatID int64, messageID int, newText string) {
	editMsg := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		newText,
	)
	bh.bot.Send(editMsg)
}

func (bh *BotHandler) deleteTgMessage(chatID int64, messageID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(
		chatID,
		messageID,
	)
	bh.bot.Send(deleteMsg)
}

func (bh *BotHandler) handleCommand(ctx context.Context, update tgbotapi.Update) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	command := tg.Command(update.Message.Command())
	args := update.Message.CommandArguments()
	switch command {
	case tg.CommandStart, tg.CommandRestart:
		err := bh.handleWelcome(ctx, tgUser)
		if err != nil {
			msg.Text = "Ошибка знакомства"
			fmt.Println(err.Error())
			bh.logger.Println(err)
			return &msg
		}

		msg.Text = "Привет, я готов запоминать задачи, начни с /add"
	case tg.CommandHelp:
		msg.Text = "Я могу /add, /list, /status и /restart."
	case tg.CommandStatus:
		msg.Text = "Я OK."
	case tg.CommandAdd:

		if args == "" {
			bh.startCommandHandling(ctx, tgUser, tg.CommandAdd)
			msg.Text = "Напиши задачу и я её запомню!"
		} else {
			_, err := bh.TaskHandler.Add(ctx, args)
			if err != nil {
				if errors.Is(err, apperrors.ErrNotFound) {
					msg.Text = err.Error()
				} else if errors.Is(err, apperrors.ErrUnAuth) {
					msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
					return &msg
				} else {
					msg.Text = "Возникла непредвиденная ошибка"
					fmt.Println(err.Error())
					bh.logger.Println(err)
				}
				break
			}
			msg.Text = fmt.Sprintf("Добавил: %v", args)
		}
	case tg.СommandTaskList:
		count, err := bh.TaskHandler.SendTaskList(ctx, bh.bot, update.Message.Chat.ID, args)
		if err != nil {
			if errors.Is(err, apperrors.ErrUnAuth) {
				msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
				return &msg
			} else {
				msg.Text = "Возникла непредвиденная ошибка"
				fmt.Println(err.Error())
				bh.logger.Println(err)
			}
			break
		}
		if count == 0 {
			msg.Text = "У меня пока нет твоих задач. Давай добавим /add"
		} else {
			msg.Text = "Вот твои задачи"
		}
	default:
		msg.Text = "Я не знаю такой команды"
	}

	return &msg
}

func (bh *BotHandler) handleMessage(ctx context.Context, update tgbotapi.Update) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	lastState, ok := bh.StateStorage.Get(ctx, tgUser.ID)
	if !ok {
		msg.Text = "Я забыл на чем мы остановились, повтори ввод команды"
		return &msg
	}

	if lastState.State != tg.StateWait {
		msg.Text = "Я уже все сделал, начни новую команду"
		return &msg
	}

	switch lastState.Command {
	case tg.CommandAdd:
		task, err := bh.TaskHandler.Add(ctx, update.Message.Text)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				msg.Text = err.Error()
			} else if errors.Is(err, apperrors.ErrUnAuth) {
				msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
				return &msg
			} else {
				msg.Text = "Возникла непредвиденная ошибка"
				fmt.Println(err.Error())
				bh.logger.Println(err)
			}
			break
		}

		msg.Text = fmt.Sprintf("Добавил: %s", task.Title)
		bh.finishCommandHandling(ctx, tgUser, tg.CommandAdd)
	default:
		msg.Text = "Для этой команды я ничего не могу сделать, начни новую"
	}

	return &msg
}

func (bh BotHandler) startCommandHandling(ctx context.Context, user *tgbotapi.User, command tg.Command) {
	bh.StateStorage.Set(ctx, user.ID, command, tg.StateWait)
}

func (bh BotHandler) finishCommandHandling(ctx context.Context, user *tgbotapi.User, command tg.Command) {
	bh.StateStorage.Set(ctx, user.ID, command, tg.StateComplete)
}

func (bh BotHandler) handleAuth(ctx context.Context, user *tgbotapi.User) context.Context {
	tokenCacheKey := "user:" + strconv.FormatInt(user.ID, 10) + ":tg-api-token-key"

	// попытка забрать из кеша
	token, _ := bh.cache.Get(ctx, tokenCacheKey).Result()
	if token != "" {
		return context.WithValue(ctx, tg.TokenContextKey, token)
	}

	// попытка сгенерировать через апи
	token, err := bh.AuthHandler.GetToken(ctx, authbot.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})
	if err != nil {
		return ctx
	}

	// вставить в кеш
	bh.cache.Set(ctx, tokenCacheKey, token, botTokenTtl)

	return context.WithValue(ctx, tg.TokenContextKey, token)
}

func (bh BotHandler) handleWelcome(ctx context.Context, user *tgbotapi.User) error {
	tokenCacheKey := "user:" + strconv.FormatInt(user.ID, 10) + ":tg-api-token-key"

	token, err := bh.AuthHandler.SignInUser(ctx, authbot.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})
	if err != nil {
		return err
	}

	bh.cache.Set(ctx, tokenCacheKey, token, botTokenTtl)

	return nil
}
