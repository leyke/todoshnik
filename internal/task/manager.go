package task

import (
	"errors"
	"sort"

	validation "todoshnik/internal/validation"
)

type TaskManager struct {
	tasks   map[int]*Task
	nextID  int
	storage TaskStorage
}

func NewTaskManager(storage TaskStorage) (*TaskManager, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, err
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

	validateError := validation.Validate(newTask)
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

func (tm *TaskManager) ListTasks(method string) []*Task {
	keys := make([]int, 0, len(tm.tasks))
	for k := range tm.tasks {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]*Task, 0, len(keys))
	for _, k := range keys {
		if method == "--pending" && tm.tasks[k].Done {
			continue
		}

		if method == "--completed" && !tm.tasks[k].Done {
			continue
		}

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
		prev := task.Done
		task.Done = true
		err := tm.storage.Save(tm.tasks)
		if err != nil {
			task.Done = prev
			return nil, err
		}

		return task, nil
	}

	return nil, errors.New("Задача не найдена")
}
