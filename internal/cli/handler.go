package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"todoshnik/internal/task"
)

func Run(tm *task.TaskManager) {
	if len(os.Args) < 2 {
		fmt.Println("Не указана команда")
		fmt.Println("Использование: go run main.go add | list | delete <Название задачи>|<ID задачи>")
		return
	}

	command := os.Args[1]
	if len(os.Args) < 3 && (command == "add" || command == "delete") {
		fmt.Printf("Команда %s требует название или ID задачи\n", command)
		return
	}
	switch command {
	case "add":
		newTitle := os.Args[2]
		task, err := tm.AddTask(newTitle)
		if err != nil {
			fmt.Printf("Ошибка при добавлении задачи: %v\n", err)
			break
		}
		fmt.Printf("Задача добавлена: ID: %d, Title: %s\n", task.ID, task.Title)
	case "list":
		method := ""
		if len(os.Args) > 2 {
			method = os.Args[2]
		}
		tasks := tm.ListTasks(method)
		if len(tasks) == 0 {
			fmt.Println("Список задач пуст")
			break
		}

		printList(tasks)
	case "done":
		taskId, err := getIntFromArgs(os.Args, 2)
		if err != nil {
			fmt.Printf("Ошибка удаления задачи: %v\n", err)
			break
		}
		_, err = tm.MarkDone(taskId)
		if err != nil {
			fmt.Printf("Ошибка пометки задачи как выполненной: %v\n", err)
		}
	case "delete":
		taskId, err := getIntFromArgs(os.Args, 2)
		if err != nil {
			fmt.Printf("Ошибка удаления задачи: %v\n", err)
			break
		}
		err = tm.DeleteTask(taskId)
		if err != nil {
			fmt.Printf("Ошибка удаления задачи: %v\n", err)
		}
	default:
		fmt.Printf("Неизвестная команда: %s\n", command)
	}
}

func getIntFromArgs(args []string, index int) (int, error) {
	if len(args) <= index {
		return 0, errors.New("Недостаточно аргументов")
	}
	return strconv.Atoi(args[index])
}

func printList(list []*task.Task) {
	for _, task := range list {
		fmt.Printf("ID: %d, Title: %s, Done: %v\n", task.ID, task.Title, task.Done)
	}
}
