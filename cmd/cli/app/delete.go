package cli

import (
	"context"

	"todoshnik/internal/infrastructure/identity"
)

func (h *Handler) delete(ctx context.Context, args []string) {
	taskId, err := getIntFromArgs(args, 2)
	if err != nil {
		printErr(err, "Ошибка удаления задачи")
		return
	}
	err = h.taskService.Delete(ctx, taskId, identity.AccessScope{IsAdmin: true})
	if err != nil {
		printErr(err, "Ошибка удаления задачи")
	}

	printInfo("Задача удалена")
}
