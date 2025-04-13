package dbUtils

import (
	"context"
	"database/sql"
	"time"
)

var QueryTimeoutDuration = 5 * time.Second

func WithTx(ctx context.Context, db *sql.DB, txFunc func(context.Context, *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = txFunc(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
