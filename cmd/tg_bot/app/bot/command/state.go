package command

type State string

const (
	StateIdle     State = "idle"
	StateWait     State = "wait"
	StateComplete State = "complete"
)
