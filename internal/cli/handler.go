package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"todoshnik/internal/domain"
	"todoshnik/internal/service"
)

type CLIHandler struct {
	service      *service.TaskService
	tokenService *service.AccessTokenService
}

func NewCLIHandler(s *service.TaskService, ts *service.AccessTokenService) *CLIHandler {
	return &CLIHandler{
		service:      s,
		tokenService: ts,
	}
}

func (cli *CLIHandler) Run() {
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
		taskId, err := getIntFromArgs(os.Args, 3)
		if err != nil {
			fmt.Printf("Ошибка получения ID задачи: %v\n", err)
			break
		}
		task, err := cli.service.AddTask(newTitle, taskId)
		if err != nil {
			fmt.Printf("Ошибка при добавлении задачи: %v\n", err)
			break
		}
		fmt.Printf("Задача добавлена: ID: %d, Title: %s\n", task.ID, task.Title)
	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		status := listCmd.String("status", "", "Фильтр по статусу: completed или pending")
		listCmd.Parse(os.Args[2:])

		tasks := cli.service.ListTasks(
			domain.TaskFilter{
				Status: domain.TaskStatus(*status),
				Scope:  domain.AccessScope{IsAdmin: true},
			})

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
		err = cli.service.MarkDone(taskId, domain.AccessScope{IsAdmin: true})
		if err != nil {
			fmt.Printf("Ошибка пометки задачи как выполненной: %v\n", err)
		}
	case "delete":
		taskId, err := getIntFromArgs(os.Args, 2)
		if err != nil {
			fmt.Printf("Ошибка удаления задачи: %v\n", err)
			break
		}
		err = cli.service.DeleteTask(taskId, domain.AccessScope{IsAdmin: true})
		if err != nil {
			fmt.Printf("Ошибка удаления задачи: %v\n", err)
		}

	case "clear-tokens":
		deleted := cli.tokenService.ClearExpiredTokens()
		fmt.Printf("Удалены токены: %v\n", deleted)
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

func printList(list []*domain.Task) {
	for _, task := range list {
		fmt.Printf("ID: %d, Title: %s, Done: %v\n", task.ID, task.Title, task.Done)
	}
}
