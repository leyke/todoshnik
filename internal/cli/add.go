package cli

import (
	"context"
	"errors"
)

func (h *Handler) add(ctx context.Context, args []string) {
	if len(args) < 3 {
		printErr(errors.New("Не указано название задачи"), "")
		return
	}

	title := args[2]

	taskID, err := getIntFromArgs(args, 3)
	if err != nil {
		printErr(err, "Ошибка ID")
		return
	}

	createdTask, err := h.taskService.Add(ctx, title, taskID)
	if err != nil {
		printErr(err, "Ошибка создания")
		return
	}

	printInfo("Задача добавлена: ")
	printTask(createdTask)
}
