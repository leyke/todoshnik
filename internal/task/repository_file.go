package task

import (
	"context"
	"sort"
	"sync"

	"todoshnik/internal/identity"
	"todoshnik/internal/storage"
)

type FileRepository struct {
	mu      sync.RWMutex
	storage storage.FileStorage[Task]
	items   map[int]*Task
	nextID  int
}

func NewFileRepository(storage storage.FileStorage[Task]) (*FileRepository, error) {
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

	return &FileRepository{
		storage: storage,
		items:   items,
		nextID:  maxID + 1,
	}, nil
}

func (repo *FileRepository) List(ctx context.Context, filter TaskFilter) ([]*Task, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	items := make([]*Task, 0)

	for _, task := range repo.items {
		if !filter.Scope.IsAdmin && task.UserID != filter.Scope.UserID {
			continue
		}

		// Фильтрация по методу
		switch filter.Status {
		case StatusPending:
			if task.Done {
				continue
			}
		case StatusCompleted:
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

	return items, nil
}

func (repo *FileRepository) GetByID(ctx context.Context, id int, scope identity.AccessScope) (*Task, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	task, ok := repo.items[id]
	if !ok || (!scope.IsAdmin && scope.UserID != task.UserID) {
		return nil, ErrNotFound
	}

	copy := *task
	return &copy, nil
}

func (repo *FileRepository) Create(ctx context.Context, task *Task) (*Task, error) {
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

	copy := *task
	return &copy, nil
}

func (repo *FileRepository) Update(ctx context.Context, task *Task) error {
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

func (repo *FileRepository) Delete(ctx context.Context, task *Task) error {
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
