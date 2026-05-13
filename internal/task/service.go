package task

import (
	"context"
	"fmt"

	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/identity"
	"todoshnik/internal/validation"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, title string, userID int) (*Task, error) {
	newTask := &Task{
		Title:  title,
		UserID: userID,
	}

	err := validateTask(newTask)
	if err != nil {
		return nil, err
	}

	task, err := s.repo.Create(ctx, newTask)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) List(ctx context.Context, filter TaskFilter) ([]*Task, error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, taskId int, title string, done bool, scope identity.AccessScope) (*Task, error) {
	task, err := s.Get(ctx, taskId, scope)
	if err != nil {
		return nil, err
	}

	prev := *task
	task.Title = title
	task.Done = done

	err = validateTask(task)
	if err != nil {
		task.Title = prev.Title
		task.Done = prev.Done
		return nil, err
	}

	err = s.repo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) Delete(ctx context.Context, taskId int, scope identity.AccessScope) error {
	task, err := s.Get(ctx, taskId, scope)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, task)
}

func (s *Service) MarkDone(ctx context.Context, taskId int, scope identity.AccessScope) (*Task, error) {
	task, err := s.Get(ctx, taskId, scope)
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

func (s *Service) Get(ctx context.Context, taskId int, scope identity.AccessScope) (*Task, error) {
	task, err := s.repo.GetByID(ctx, taskId, scope)
	return task, err
}

func validateTask(task *Task) error {
	ve := validation.Validate(task)
	if ve != nil {
		fmt.Println(ve)
		return apperrors.NewValidationErrorFromValidator(ve)
	}
	return nil
}
