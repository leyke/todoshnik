package token

import "time"

// Конфиг есть, но почему-то ты не возмользовался?
type Config struct {
	Secret string
	Ttl    time.Duration
}
