package service

import (
	"errors"
	"sort"

	"todoshnik/internal/domain"
	"todoshnik/internal/storage"

	"todoshnik/internal/validation"
)

type TaskService struct {
	tasks   map[int]*domain.Task
	nextID  int
	storage storage.TaskStorage
}

func NewTaskService(storage storage.TaskStorage) (*TaskService, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, err
	}

	s := &TaskService{
		tasks:   tasks,
		storage: storage,
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	s.nextID = maxID + 1

	return s, nil
}

func (s *TaskService) AddTask(title string) (*domain.Task, error) {
	newTask := &domain.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}

	validateError := validation.Validate(newTask)
	if validateError != nil {
		return nil, validateError
	}

	s.tasks[s.nextID] = newTask
	s.nextID++

	err := s.storage.Save(s.tasks)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func (s *TaskService) ListTasks(method string) []*domain.Task {
	keys := make([]int, 0, len(s.tasks))
	for k := range s.tasks {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]*domain.Task, 0, len(keys))
	for _, k := range keys {
		if method == "pending" && s.tasks[k].Done {
			continue
		}

		if method == "completed" && !s.tasks[k].Done {
			continue
		}

		result = append(result, s.tasks[k])
	}
	return result
}

func (s *TaskService) DeleteTask(taskId int) error {
	_, ok := s.tasks[taskId]
	if ok {
		delete(s.tasks, taskId)

		return s.storage.Save(s.tasks)
	}

	return errors.New("Задача не найдена")
}

func (s *TaskService) MarkDone(taskId int) (*domain.Task, error) {
	_, ok := s.tasks[taskId]
	if ok {
		task := s.tasks[taskId]
		prev := task.Done
		task.Done = true
		err := s.storage.Save(s.tasks)
		if err != nil {
			task.Done = prev
			return nil, err
		}

		return task, nil
	}

	return nil, errors.New("Задача не найдена")
}
