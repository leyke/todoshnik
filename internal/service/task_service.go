package service

import (
	"fmt"
	"os"

	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/storage"
	"todoshnik/internal/validation"

	repo "todoshnik/internal/repository/task"

	"github.com/go-playground/validator/v10"
)

type TaskService struct {
	repo repo.TaskRepositoryInreface
}

func NewTaskService() (*TaskService, error) {
	storagePath := os.Getenv("tmp_dir") + "/tasks.json"
	storage := storage.NewFileStorage[domain.Task](storagePath)

	repo, err := repo.NewTaskFileRepository(storage)
	if err != nil {
		return nil, err
	}

	s := &TaskService{
		repo: repo,
	}

	return s, nil
}

func (s *TaskService) AddTask(title string, userID int) (*domain.Task, error) {
	newTask := &domain.Task{
		Title:  title,
		UserID: userID,
	}

	ve := validation.Validate(newTask)
	if ve != nil {
		fmt.Println(ve)
		return nil, apperrors.NewValidationErrorFromValidator(ve.(validator.ValidationErrors))
	}

	task, err := s.repo.Create(newTask)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) ListTasks(filter domain.TaskFilter) []*domain.Task {
	return s.repo.List(filter)
}

func (s *TaskService) UpdateTask(taskId int, title string, done bool, scope domain.AccessScope) (*domain.Task, error) {
	task, err := s.GetTask(taskId, scope)
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

	err = s.repo.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(taskId int, scope domain.AccessScope) error {
	task, err := s.GetTask(taskId, scope)
	if err != nil {
		return err
	}

	return s.repo.Delete(task)
}

func (s *TaskService) MarkDone(taskId int, scope domain.AccessScope) error {
	task, err := s.GetTask(taskId, scope)
	if err != nil {
		return err
	}

	task.Done = !task.Done
	err = s.repo.Update(task)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) GetTask(taskId int, scope domain.AccessScope) (*domain.Task, error) {
	task, err := s.repo.GetByID(taskId, scope)
	return task, err
}
