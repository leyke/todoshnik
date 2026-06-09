package cli

import (
	"context"
	"flag"

	"todoshnik/internal/identity"
	"todoshnik/internal/task"
)

func (h *Handler) list(ctx context.Context, args []string) {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	status := listCmd.String("status", "", "Фильтр по статусу: completed или pending")
	listCmd.Parse(args[2:])

	tasks, err := h.taskService.List(
		ctx,
		task.TaskFilter{
			Status: task.Status(*status),
			Scope:  identity.AccessScope{IsAdmin: true},
		})

	if err != nil {
		printErr(err, "Ошибка получения списка задач")
		return
	}

	if len(tasks) == 0 {
		printInfo("Список задач пуст")
		return
	}

	printList(tasks)
}
