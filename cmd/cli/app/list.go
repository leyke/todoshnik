package cli

import (
	"context"
	"flag"

	"todoshnik/internal/domains/task"
	"todoshnik/internal/infrastructure/identity"
)

func (h *Handler) list(ctx context.Context, args []string) {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	status := listCmd.String("status", "", "Фильтр по статусу: completed или pending")
	err := listCmd.Parse(args[2:])
	if err != nil {
		printErr(err, "ошибка получения параметров команды")
		return
	}

	tasks, err := h.taskService.List(
		ctx,
		task.TaskFilter{
			Status: task.Status(*status),
			Scope:  identity.AccessScope{IsAdmin: true},
		})

	if err != nil {
		printErr(err, "ошибка получения списка задач")
		return
	}

	if len(tasks) == 0 {
		printInfo("Список задач пуст")
		return
	}

	printList(tasks)
}
