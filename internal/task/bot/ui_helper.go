package bot

import (
	"fmt"

	"todoshnik/internal/task"
)

const (
	emojiInProgress string = "⬜"
	emojiIsDone     string = "✅"
	emojiDelete     string = "🗑️"
)

func GetStatusButtonText(task task.Task) string {
	if task.Done {
		return emojiInProgress + " В процессе"
	}
	return emojiIsDone + " Готово"
}

func GetDeleteButtonText() string {
	return emojiDelete + " Забыть"
}

func GetTaskRowText(task task.Task) string {
	emoji := emojiInProgress
	if task.Done {
		emoji = emojiIsDone
	}
	return fmt.Sprintf("%s %s", emoji, task.Title)
}
