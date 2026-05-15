package session

import "todoshnik/internal/bot/tg"

type Session struct {
	tg.Command
	tg.State
}
