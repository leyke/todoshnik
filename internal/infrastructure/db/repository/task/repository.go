package task

import (
	"context"
	"database/sql"
	"errors"

	"todoshnik/internal/domains/task"
	"todoshnik/internal/infrastructure/identity"

	taskerrors "todoshnik/internal/domains/task/errors"
	db "todoshnik/internal/infrastructure/db/transaction"

	"github.com/Masterminds/squirrel"
)

type DBRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *DBRepository {
	return &DBRepository{
		db: db,
	}
}

func (repo *DBRepository) List(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error) {
	builder := squirrel.
		Select("id", "title", "done", "user_id").
		PlaceholderFormat(squirrel.Dollar).
		From("tasks")

	if !filter.Scope.IsAdmin {
		builder = builder.Where(squirrel.Eq{
			"user_id": filter.Scope.UserID,
		})
	}

	switch filter.Status {
	case task.StatusPending:
		builder = builder.Where(squirrel.Eq{
			"done": false,
		})

	case task.StatusCompleted:
		builder = builder.Where(squirrel.Eq{
			"done": true,
		})
	}

	builder = builder.
		OrderBy("done DESC", "id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	rows, err := executor.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var items []*task.Task
	for rows.Next() {
		item := &task.Task{}

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Done,
			&item.UserID,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repo *DBRepository) GetByID(ctx context.Context, id int, scope identity.AccessScope) (*task.Task, error) {
	builder := squirrel.
		Select("id", "title", "done", "user_id").
		PlaceholderFormat(squirrel.Dollar).
		From("tasks").
		Where(squirrel.Eq{"id": id})

	if !scope.IsAdmin {
		builder = builder.Where(squirrel.Eq{"user_id": scope.UserID})
	}

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	item := &task.Task{}

	err = executor.
		QueryRowContext(ctx, query, args...).
		Scan(
			&item.ID,
			&item.Title,
			&item.Done,
			&item.UserID,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, taskerrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (repo *DBRepository) Create(ctx context.Context, t *task.Task) (*task.Task, error) {
	query, args, err := squirrel.
		Insert("tasks").
		Columns("title", "done", "user_id").
		Values(t.Title, t.Done, t.UserID).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	err = executor.
		QueryRowContext(
			ctx,
			query,
			args...,
		).
		Scan(&t.ID)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (repo *DBRepository) Update(ctx context.Context, t *task.Task) error {
	query, args, err := squirrel.
		Update("tasks").
		PlaceholderFormat(squirrel.Dollar).
		Set("title", t.Title).
		Set("done", t.Done).
		Set("user_id", t.UserID).
		Where(squirrel.Eq{"id": t.ID}).
		ToSql()

	if err != nil {
		return err
	}

	// Выполняем запрос
	executor := db.ExecutorFromContext(ctx, repo.db)

	result, err := executor.
		ExecContext(
			ctx,
			query,
			args...,
		)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return taskerrors.ErrNotFound
	}

	return nil
}

func (repo *DBRepository) Delete(ctx context.Context, t *task.Task) error {
	query, args, err := squirrel.
		Delete("tasks").
		Where(squirrel.Eq{"id": t.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	// Выполняем запрос
	executor := db.ExecutorFromContext(ctx, repo.db)

	_, err = executor.
		ExecContext(
			ctx,
			query,
			args...,
		)

	if err != nil {
		return err
	}

	return nil
}
