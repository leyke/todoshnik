package task

import (
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

func (repo *TaskDbRepository) List(filter domain.TaskFilter) []*domain.Task {
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
	repo.db.Order("done desc, id").Find(result)

	return result
}

func (repo *TaskDbRepository) GetByID(id int, scope domain.AccessScope) (*domain.Task, error) {
	var item *domain.Task
	result := repo.db.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *TaskDbRepository) Create(task *domain.Task) (*domain.Task, error) {
	result := repo.db.Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (repo *TaskDbRepository) Update(task *domain.Task) error {
	result := repo.db.Save(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *TaskDbRepository) Delete(task *domain.Task) error {
	result := repo.db.Delete(task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
