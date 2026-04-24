package service

import (
	"fmt"
	"os"
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
	storage *storage.FileStorage[domain.Task]
	mu      sync.Mutex
}

func NewTaskService() (*TaskService, error) {
	storagePath := os.Getenv("tmp_dir") + "/tasks.json"
	storage := storage.NewFileStorage[domain.Task](storagePath)

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

func (s *TaskService) AddTask(title string, userID *int) (*domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newTask := &domain.Task{
		ID:     s.nextID,
		Title:  title,
		Done:   false,
		UserID: *userID,
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

func (s *TaskService) ListTasks(method string, userId *int) []*domain.Task {
	tasks := make([]*domain.Task, 0)

	for _, task := range s.tasks {
		if userId != nil && task.UserID != *userId {
			continue
		}

		// Фильтрация по методу
		switch method {
		case "pending":
			if task.Done {
				continue
			}
		case "completed":
			if !task.Done {
				continue
			}
		}

		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Done != tasks[j].Done {
			return tasks[i].Done
		}
		return tasks[i].ID < tasks[j].ID
	})

	return tasks
}

func (s *TaskService) UpdateTask(taskId int, title string, done bool, userId int) (*domain.Task, error) {
	task, err := s.GetTask(taskId, &userId)
	if err != nil {
		return nil, err
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

	err = s.storage.Save(s.tasks)
	if err != nil {
		s.tasks[taskId] = prev
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(taskId int, userId *int) error {
	_, err := s.GetTask(taskId, userId)
	if err != nil {
		return err
	}
	delete(s.tasks, taskId)

	return s.storage.Save(s.tasks)
}

func (s *TaskService) MarkDone(taskId int, userId *int) (*domain.Task, error) {
	task, err := s.GetTask(taskId, userId)
	if err != nil {
		return nil, err
	}

	prev := task.Done
	task.Done = !prev
	err = s.storage.Save(s.tasks)
	if err != nil {
		task.Done = prev
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(taskId int, userId *int) (*domain.Task, error) {
	task, ok := s.tasks[taskId]
	if !ok || (userId != nil && task.UserID != *userId) {
		return nil, apperrors.ErrNotFound
	}

	return task, nil
}
