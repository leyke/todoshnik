package token

import (
	"context"
	"database/sql"
	"errors"
	"time"

	apptoken "todoshnik/internal/domains/token"
	tokenerror "todoshnik/internal/domains/token/errors"
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

func (repo *DBRepository) GetAllByUserID(ctx context.Context, userID int) ([]*apptoken.Token, error) {
	builder := squirrel.
		Select("id", "user_id", "hash", "device", "expires_at").
		PlaceholderFormat(squirrel.Dollar).
		From("tokens").
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

	var items []*apptoken.Token
	for rows.Next() {
		item := &apptoken.Token{}

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Hash,
			&item.Device,
			&item.ExpiresAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repo *DBRepository) GetByHash(ctx context.Context, hash string) (*apptoken.Token, error) {
	builder := squirrel.
		Select("id", "user_id", "hash", "device", "expires_at").
		PlaceholderFormat(squirrel.Dollar).
		From("tokens").
		Where(squirrel.Eq{"hash": hash})

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	executor := db.ExecutorFromContext(ctx, repo.db)

	item := &apptoken.Token{}

	err = executor.
		QueryRowContext(ctx, query, args...).
		Scan(
			&item.ID,
			&item.UserID,
			&item.Hash,
			&item.Device,
			&item.ExpiresAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, tokenerror.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (repo *DBRepository) GetExpiredTokens(ctx context.Context, before time.Time) ([]*apptoken.Token, error) {
	builder := squirrel.
		Select("id", "user_id", "hash", "device", "expires_at").
		From("tokens").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Lt{"expires_at": before}).
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

	var items []*apptoken.Token
	for rows.Next() {
		item := &apptoken.Token{}

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Hash,
			&item.Device,
			&item.ExpiresAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repo *DBRepository) Create(ctx context.Context, t *apptoken.Token) (*apptoken.Token, error) {
	query, args, err := squirrel.
		Insert("tokens").
		Columns("user_id", "hash", "device", "expires_at").
		Values(t.UserID, t.Hash, t.Device, t.ExpiresAt).
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

	if db.IsForeignKeyViolation(err) {
		return nil, tokenerror.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (repo *DBRepository) Delete(ctx context.Context, t *apptoken.Token) error {
	query, args, err := squirrel.
		Delete("tokens").
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
