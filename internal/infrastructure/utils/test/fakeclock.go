package clock

import "time"

type FakeClock struct {
	Time time.Time
}

func (c *FakeClock) Now() time.Time {
	return c.Time
}
