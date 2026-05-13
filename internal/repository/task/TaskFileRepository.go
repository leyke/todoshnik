package task

import (
	"context"
	"sort"
	"sync"
	"todoshnik/internal/domain"
	apperror "todoshnik/internal/errors"
	"todoshnik/internal/storage"
)

type TaskFileRepository struct {
	mu      sync.RWMutex
	storage storage.FileStorage[domain.Task]
	items   map[int]*domain.Task
	nextID  int
}

func NewTaskFileRepository(storage storage.FileStorage[domain.Task]) (*TaskFileRepository, error) {
	items, err := storage.Load()
	if err != nil {
		return nil, err
	}

	maxID := 0
	for _, task := range items {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return &TaskFileRepository{
		storage: storage,
		items:   items,
		nextID:  maxID + 1,
	}, nil
}

func (repo *TaskFileRepository) List(ctx context.Context, filter domain.TaskFilter) []*domain.Task {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	items := make([]*domain.Task, 0)

	for _, task := range repo.items {
		if !filter.Scope.IsAdmin && task.UserID != filter.Scope.UserID {
			continue
		}

		// Фильтрация по методу
		switch filter.Status {
		case domain.StatusPending:
			if task.Done {
				continue
			}
		case domain.StatusCompleted:
			if !task.Done {
				continue
			}
		}

		items = append(items, task)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Done != items[j].Done {
			return items[i].Done
		}
		return items[i].ID < items[j].ID
	})

	return items
}

func (repo *TaskFileRepository) GetByID(ctx context.Context, id int, scope domain.AccessScope) (*domain.Task, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	task, ok := repo.items[id]
	if !ok || (!scope.IsAdmin && scope.UserID != task.UserID) {
		return nil, apperror.ErrNotFound
	}

	return task, nil
}

func (repo *TaskFileRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	task.ID = repo.nextID
	repo.nextID++

	repo.items[task.ID] = task

	err := repo.storage.Save(repo.items)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (repo *TaskFileRepository) Update(ctx context.Context, task *domain.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	prev := repo.items[task.ID]
	repo.items[task.ID] = task

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[task.ID] = prev
		return err
	}

	return nil
}

func (repo *TaskFileRepository) Delete(ctx context.Context, task *domain.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	prev := task
	delete(repo.items, task.ID)

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[task.ID] = prev
		return err
	}

	return nil
}
