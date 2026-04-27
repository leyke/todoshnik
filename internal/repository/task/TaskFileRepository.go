package task

import (
	"sort"
	"sync"
	"todoshnik/internal/domain"
	apperror "todoshnik/internal/errors"
	"todoshnik/internal/storage"
)

type TaskFileRepository struct {
	mu      sync.RWMutex
	storage storage.FileStorage[domain.Task]
	tasks   map[int]*domain.Task
	nextID  int
}

func NewTaskFileRepository(storage storage.FileStorage[domain.Task]) (*TaskFileRepository, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, err
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return &TaskFileRepository{
		storage: storage,
		tasks:   tasks,
		nextID:  maxID + 1,
	}, nil
}

func (repo *TaskFileRepository) List(filter domain.TaskFilter) []*domain.Task {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	tasks := make([]*domain.Task, 0)

	for _, task := range repo.tasks {
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

func (repo *TaskFileRepository) GetByID(id int, scope domain.AccessScope) (*domain.Task, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	task, ok := repo.tasks[id]
	if !ok || (!scope.IsAdmin && scope.UserID != task.UserID) {
		return nil, apperror.ErrNotFound
	}

	return task, nil
}

func (repo *TaskFileRepository) Create(task *domain.Task) (*domain.Task, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	task.ID = repo.nextID
	repo.nextID++

	repo.tasks[task.ID] = task

	err := repo.storage.Save(repo.tasks)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (repo *TaskFileRepository) Update(task *domain.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	prev := repo.tasks[task.ID]
	repo.tasks[task.ID] = task

	err := repo.storage.Save(repo.tasks)
	if err != nil {
		repo.tasks[task.ID] = prev
		return err
	}

	return nil
}

func (repo *TaskFileRepository) Delete(task *domain.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	prev := task
	delete(repo.tasks, task.ID)

	err := repo.storage.Save(repo.tasks)
	if err != nil {
		repo.tasks[task.ID] = prev
		return err
	}

	return nil
}
