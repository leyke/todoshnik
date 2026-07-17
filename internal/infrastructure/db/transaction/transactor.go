package transaction

import (
	"context"
	"database/sql"
)

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

func (t *Transactor) WithinTransaction(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	ctx = injectTx(ctx, tx)

	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
