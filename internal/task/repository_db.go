package task

import (
	"context"
	"todoshnik/internal/identity"

	"gorm.io/gorm"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDbRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{
		db: db,
	}
}

func (repo *DBRepository) List(ctx context.Context, filter TaskFilter) ([]*Task, error) {
	var items []*Task

	where := make(map[string]any)

	if !filter.Scope.IsAdmin {
		where["user_id"] = filter.Scope.UserID
	}

	// Фильтрация по методу
	switch filter.Status {
	case StatusPending:
		where["done"] = false
	case StatusCompleted:
		where["done"] = true
	}
	result := repo.db.WithContext(ctx).Order("done desc, id").Find(&items, where)
	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

func (repo *DBRepository) GetByID(ctx context.Context, id int, scope identity.AccessScope) (*Task, error) {
	var item *Task

	query := repo.db.WithContext(ctx)

	if !scope.IsAdmin {
		query = query.Where("user_id = ?", scope.UserID)
	}

	result := query.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *DBRepository) Create(ctx context.Context, task *Task) (*Task, error) {
	result := repo.db.WithContext(ctx).Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (repo *DBRepository) Update(ctx context.Context, task *Task) error {
	result := repo.db.WithContext(ctx).Save(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DBRepository) Delete(ctx context.Context, task *Task) error {
	result := repo.db.WithContext(ctx).Delete(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
