package clock

import "time"

type RealClock struct{}

func New() *RealClock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}
