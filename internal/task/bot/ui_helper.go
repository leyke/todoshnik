package bot

import (
	"fmt"
	"todoshnik/internal/constants"
	"todoshnik/internal/task"
)

func getStatusButtonText(task task.Task) string {
	if task.Done {
		return constants.EmojiInProgress + " В процессе"
	}
	return constants.EmojiIsDone + " Готово"
}

func getDeleteButtonText() string {
	return constants.EmojiDelete + " Забыть"
}

func getTaskRowText(task task.Task) string {
	emoji := constants.EmojiInProgress
	if task.Done {
		emoji = constants.EmojiIsDone
	}
	return fmt.Sprintf("%s %s", emoji, task.Title)
}
