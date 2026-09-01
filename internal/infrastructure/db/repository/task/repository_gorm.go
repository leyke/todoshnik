package task

import (
	"context"
	"errors"

	"todoshnik/internal/domains/task"
	"todoshnik/internal/infrastructure/identity"

	taskerrors "todoshnik/internal/domains/task/errors"

	"gorm.io/gorm"
)

// deprecated оставил для примера работы с gorm
type GormDBRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormDBRepository {
	return &GormDBRepository{
		db: db,
	}
}

func (repo *GormDBRepository) List(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error) {
	var items []*task.Task

	where := make(map[string]any)

	if !filter.Scope.IsAdmin {
		where["user_id"] = filter.Scope.UserID
	}

	// Фильтрация по методу
	switch filter.Status {
	case task.StatusPending:
		where["done"] = false
	case task.StatusCompleted:
		where["done"] = true
	}
	result := repo.db.WithContext(ctx).Order("done desc, id").Find(&items, where)
	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

func (repo *GormDBRepository) GetByID(ctx context.Context, id int, scope identity.AccessScope) (*task.Task, error) {
	var item *task.Task

	query := repo.db.WithContext(ctx)

	if !scope.IsAdmin {
		query = query.Where("user_id = ?", scope.UserID)
	}

	result := query.First(&item, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, taskerrors.ErrNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *GormDBRepository) Create(ctx context.Context, task *task.Task) (*task.Task, error) {
	result := repo.db.WithContext(ctx).Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (repo *GormDBRepository) Update(ctx context.Context, task *task.Task) error {
	result := repo.db.WithContext(ctx).Save(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *GormDBRepository) Delete(ctx context.Context, task *task.Task) error {
	result := repo.db.WithContext(ctx).Delete(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
