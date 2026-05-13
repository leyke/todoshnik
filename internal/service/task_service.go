package service

import (
	"context"
	"fmt"

	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/validation"

	repository "todoshnik/internal/repository/task"

	"github.com/go-playground/validator/v10"
)

type TaskService struct {
	repo repository.TaskRepositoryInterface
}

func NewTaskService(repo repository.TaskRepositoryInterface) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) AddTask(ctx context.Context, title string, userID int) (*domain.Task, error) {
	newTask := &domain.Task{
		Title:  title,
		UserID: userID,
	}

	ve := validation.Validate(newTask)
	if ve != nil {
		fmt.Println(ve)
		return nil, apperrors.NewValidationErrorFromValidator(ve.(validator.ValidationErrors))
	}

	task, err := s.repo.Create(ctx, newTask)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, filter domain.TaskFilter) []*domain.Task {
	return s.repo.List(ctx, filter)
}

func (s *TaskService) UpdateTask(ctx context.Context, taskId int, title string, done bool, scope domain.AccessScope) (*domain.Task, error) {
	task, err := s.GetTask(ctx, taskId, scope)
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

	err = s.repo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskId int, scope domain.AccessScope) error {
	task, err := s.GetTask(ctx, taskId, scope)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, task)
}

func (s *TaskService) MarkDone(ctx context.Context, taskId int, scope domain.AccessScope) (*domain.Task, error) {
	task, err := s.GetTask(ctx, taskId, scope)
	if err != nil {
		return nil, err
	}

	task.Done = !task.Done
	err = s.repo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskId int, scope domain.AccessScope) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, taskId, scope)
	return task, err
}
