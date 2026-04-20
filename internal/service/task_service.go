package service

import (
	"fmt"
	"sort"
	"sync"

	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/storage"

	"todoshnik/internal/validation"

	"github.com/go-playground/validator/v10"
)

type TaskService struct {
	tasks   map[int]*domain.Task
	nextID  int
	storage storage.TaskStorage
	mu      sync.Mutex
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
	s.mu.Lock()
	defer s.mu.Unlock()

	newTask := &domain.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}

	ve := validation.Validate(newTask)
	if ve != nil {
		fmt.Println(ve)
		return nil, apperrors.NewValidationErrorFromValidator(ve.(validator.ValidationErrors))
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

func (s *TaskService) UpdateTask(taskId int, title string, done bool) (*domain.Task, error) {
	task, ok := s.tasks[taskId]
	if !ok {
		return nil, apperrors.ErrNotFound
	}

	prev := task
	task.Title = title
	task.Done = done
	validateError := validation.Validate(task)
	if validateError != nil {
		task.Title = prev.Title
		task.Done = prev.Done
		return nil, apperrors.NewValidationErrorFromValidator(validateError.(validator.ValidationErrors))
	}

	err := s.storage.Save(s.tasks)
	if err != nil {
		s.tasks[taskId] = prev
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(taskId int) error {
	_, ok := s.tasks[taskId]
	if !ok {
		return apperrors.ErrNotFound
	}

	delete(s.tasks, taskId)

	return s.storage.Save(s.tasks)
}

func (s *TaskService) MarkDone(taskId int) (*domain.Task, error) {
	_, ok := s.tasks[taskId]
	if !ok {
		return nil, apperrors.ErrNotFound
	}

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

func (s *TaskService) GetTask(taskId int) (*domain.Task, error) {
	task, ok := s.tasks[taskId]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return task, nil
}
