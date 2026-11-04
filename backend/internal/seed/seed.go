package seed

import (
	"context"
	"fmt"
	"io/fs"
	"sort"

	"github.com/cephalic-labs/nephrosa/db/seed"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Apply(ctx context.Context, pool *pgxpool.Pool) error {
	files, err := fs.Glob(seed.Files, "*.sql")
	if err != nil {
		return fmt.Errorf("list embedded seed files: %w", err)
	}
	sort.Strings(files)

	for _, name := range files {
		sql, err := seed.Files.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read seed file %s: %w", name, err)
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin seed %s: %w", name, err)
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("apply seed %s: %w", name, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit seed %s: %w", name, err)
		}
		fmt.Printf("applied fictional seed %s\n", name)
	}
	return nil
}
