package task

import (
	"context"
	"todoshnik/internal/domain"

	"gorm.io/gorm"
)

type TaskDbRepository struct {
	db *gorm.DB
}

func NewTaskDbRepository(db *gorm.DB) *TaskDbRepository {
	return &TaskDbRepository{
		db: db,
	}
}

func (repo *TaskDbRepository) List(ctx context.Context, filter domain.TaskFilter) []*domain.Task {
	var result []*domain.Task

	where := make(map[string]any)

	if !filter.Scope.IsAdmin {
		where["user_id"] = filter.Scope.UserID
	}

	// Фильтрация по методу
	switch filter.Status {
	case domain.StatusPending:
		where["done"] = false
	case domain.StatusCompleted:
		where["done"] = true
	}
	repo.db.WithContext(ctx).Order("done desc, id").Find(&result, where)

	return result
}

func (repo *TaskDbRepository) GetByID(ctx context.Context, id int, scope domain.AccessScope) (*domain.Task, error) {
	var item *domain.Task
	result := repo.db.WithContext(ctx).First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *TaskDbRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	result := repo.db.WithContext(ctx).Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (repo *TaskDbRepository) Update(ctx context.Context, task *domain.Task) error {
	result := repo.db.WithContext(ctx).Save(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *TaskDbRepository) Delete(ctx context.Context, task *domain.Task) error {
	result := repo.db.WithContext(ctx).Delete(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
