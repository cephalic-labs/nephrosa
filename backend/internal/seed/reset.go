package seed

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Reset requires the migration role; the caller must check environment and explicit confirmation.
func Reset(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(ctx, `LOCK TABLE patients,episodes,call_attempts,users,sessions IN ACCESS EXCLUSIVE MODE`)
	if err != nil {
		return err
	}
	var unsafe bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM patients WHERE external_reference NOT LIKE 'DEMO-%')`).Scan(&unsafe)
	if err != nil {
		return err
	}
	if unsafe {
		return errors.New("reset refused: non-demo patient references exist")
	}
	_, err = tx.Exec(ctx, `TRUNCATE sessions,call_attempts,audit_events,clinical_reviews,reports,test_completions,tasks,barriers,contact_attempts,test_orders,episodes,consents,patients,users RESTART IDENTITY`)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
