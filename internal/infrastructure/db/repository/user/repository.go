package user

import (
	"context"
	"database/sql"
	"errors"

	appuser "todoshnik/internal/domains/user"
	usererrors "todoshnik/internal/domains/user/errors"
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

func (repo *DBRepository) List(ctx context.Context) ([]*appuser.User, error) {
	builder := squirrel.
		Select("id", "name", "telegram_id", "login").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("id")

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

	var items []*appuser.User

	for rows.Next() {
		item := &appuser.User{}
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TelegramID,
			&item.Login,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repo *DBRepository) GetByID(ctx context.Context, id int) (*appuser.User, error) {
	builder := squirrel.
		Select("id", "name", "telegram_id", "login").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id})

	query, args, err := builder.
		ToSql()

	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	item := &appuser.User{}

	err = executor.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.Name,
		&item.TelegramID,
		&item.Login,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, usererrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (repo *DBRepository) GetByLogin(ctx context.Context, login string) (*appuser.User, error) {
	builder := squirrel.
		Select("id", "name", "telegram_id", "login").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"login": login})

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	item := &appuser.User{}

	err = executor.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.Name,
		&item.TelegramID,
		&item.Login,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, usererrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (repo *DBRepository) GetByTgId(ctx context.Context, userTgId int64) (*appuser.User, error) {
	builder := squirrel.
		Select("id", "name", "telegram_id", "login").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"telegram_id": userTgId})

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	item := &appuser.User{}

	err = executor.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.Name,
		&item.TelegramID,
		&item.Login,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, usererrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (repo *DBRepository) Create(ctx context.Context, u *appuser.User) (*appuser.User, error) {
	query, args, err := squirrel.
		Insert("users").
		PlaceholderFormat(squirrel.Dollar).
		Columns("name", "telegram_id", "login", "password_hash").
		Values(u.Name, u.TelegramID, u.Login, u.PasswordHash).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	err = executor.
		QueryRowContext(ctx, query, args...).
		Scan(&u.ID)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *DBRepository) Update(ctx context.Context, u *appuser.User) error {
	query, args, err := squirrel.
		Update("users").
		PlaceholderFormat(squirrel.Dollar).
		Set("name", u.Name).
		Set("telegram_id", u.TelegramID).
		Set("login", u.Login).
		Set("passwordHash", u.PasswordHash).
		Where(squirrel.Eq{"id": u.ID}).
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

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return usererrors.ErrNotFound
	}

	return nil
}

func (repo *DBRepository) Delete(ctx context.Context, u *appuser.User) error {
	query, args, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{"id": u.ID}).
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
