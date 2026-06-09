package cli

import (
	"context"

	"todoshnik/internal/identity"
)

func (h *Handler) done(ctx context.Context, args []string) {
	taskId, err := getIntFromArgs(args, 2)
	if err != nil {
		printErr(err, "Ошибка ID")
		return
	}
	task, err := h.taskService.MarkDone(ctx, taskId, identity.AccessScope{IsAdmin: true})
	if err != nil {
		printErr(err, "Ошибка пометки задачи как выполненной")

	}

	printInfo("Задача помечена как выполненная: ")
	printTask(task)
}
