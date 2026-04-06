package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title" validate:"required,min=3"`
	Done  bool   `json:"done,omitempty"`
}

type TaskManager struct {
	tasks   map[int]*Task
	nextID  int
	storage TaskStorage
}

type TaskStorage interface {
	Save(tasks map[int]*Task) error
	Load() (map[int]*Task, error)
}

type FileStorage struct {
	filename string
}

func (fs *FileStorage) Save(tasks map[int]*Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filename, data, 0644)
}

func (fs *FileStorage) Load() (map[int]*Task, error) {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[int]*Task), nil
		}
		return nil, err
	}

	tasks := make(map[int]*Task)
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func NewTaskManager(storage TaskStorage) (*TaskManager, error) {
	tasks, error := storage.Load()
	if error != nil {
		return nil, error
	}

	tm := &TaskManager{
		tasks:   tasks,
		storage: storage,
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	tm.nextID = maxID + 1

	return tm, nil
}

func (tm *TaskManager) AddTask(title string) (*Task, error) {
	newTask := &Task{
		ID:    tm.nextID,
		Title: title,
		Done:  false,
	}

	validate := validator.New()
	validateError := validate.Struct(newTask)
	if validateError != nil {
		return nil, validateError
	}

	tm.tasks[tm.nextID] = newTask
	tm.nextID++

	err := tm.storage.Save(tm.tasks)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func (tm *TaskManager) ListTasks() []*Task {
	keys := make([]int, 0, len(tm.tasks))
	for k := range tm.tasks {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]*Task, 0, len(keys))
	for _, k := range keys {
		result = append(result, tm.tasks[k])
	}
	return result
}

func (tm *TaskManager) DeleteTask(taskId int) error {
	_, ok := tm.tasks[taskId]
	if ok {
		delete(tm.tasks, taskId)

		return tm.storage.Save(tm.tasks)
	}

	return errors.New("Задача не найдена")
}

func (tm *TaskManager) MarkDone(taskId int) (*Task, error) {
	_, ok := tm.tasks[taskId]
	if ok {
		task := tm.tasks[taskId]
		task.Done = true
		err := tm.storage.Save(tm.tasks)
		if err != nil {
			task.Done = false
			return nil, err
		}
		return task, nil
	}

	return nil, errors.New("Задача не найдена")
}

func main() {
	storage := &FileStorage{filename: "tmp/tasks.json"}
	tm, err := NewTaskManager(storage)
	if err != nil {
		fmt.Printf("Ошибка инициализации менеджера задач: %v\n", err)
		return
	}

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
		task, error := tm.AddTask(newTitle)
		if error != nil {
			fmt.Printf("Ошибка при добавлении задачи: %v\n", error)
			break
		}
		fmt.Printf("Задача добавлена: ID: %d, Title: %s\n", task.ID, task.Title)
	case "list":
		tasks := tm.ListTasks()
		if len(tasks) == 0 {
			fmt.Println("Список задач пуст")
			break
		}

		printList(tasks)
	case "done":
		taskId, err := GetIntFromArgs(os.Args, 2)
		if err != nil {
			fmt.Printf("Ошибка удаления задачи: %v\n", err)
			break
		}
		_, err = tm.MarkDone(taskId)
		if err != nil {
			fmt.Printf("Ошибка пометки задачи как выполненной: %v\n", err)
		}
	case "delete":
		taskId, err := GetIntFromArgs(os.Args, 2)
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

func GetIntFromArgs(args []string, index int) (int, error) {
	if len(args) <= index {
		return 0, errors.New("Недостаточно аргументов")
	}
	return strconv.Atoi(args[index])
}

func printList(list []*Task) {
	for _, task := range list {
		fmt.Printf("ID: %d, Title: %s, Done: %v\n", task.ID, task.Title, task.Done)
	}
}
