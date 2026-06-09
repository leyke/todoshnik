package cli

import (
	"fmt"

	"todoshnik/internal/task"
)

func printList(list []*task.Task) {
	for _, task := range list {
		printTask(task)
	}
}

func printTask(task *task.Task) {
	fmt.Printf("ID: %d, Title: %s, Done: %v\n", task.ID, task.Title, task.Done)
}

func printInfo(msg string) {
	fmt.Println(msg)
}

func printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("  add <название> <ID> - Добавить задачу")
	fmt.Println("  list [-status <статус>] - Показать список задач")
	fmt.Println("  done <ID> - Пометить задачу как выполненную")
	fmt.Println("  delete <ID> - Удалить задачу")
	fmt.Println("  clear-tokens - Очистить истекшие токены")
}

func printErr(err error, msg string) {
	if msg == "" {
		msg = "Ошибка"
	}

	fmt.Printf("%s: %v\n", msg, err)
}
