package callback

import "todoshnik/cmd/tg_bot/app/bot/command"

type CallbackQuery struct {
	Command command.Name      `json:"command"`
	Payload map[string]string `json:"payload"`
}
