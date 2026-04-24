package tg

type CallbackQuery struct {
	Command Command           `json:"Command"`
	Payload map[string]string `json:"Payload"`
}
