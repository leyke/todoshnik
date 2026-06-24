package token

import "time"

type Config struct {
	Secret string
	Ttl    time.Duration
}
