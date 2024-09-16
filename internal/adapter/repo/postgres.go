package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const defaultTimeout = time.Minute

var (
	errFilterIsEmpty = errors.New("filter is empty")
)

var defaultTables = map[string]string{
	"uuid": `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
}

type PostgresRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func (r *PostgresRepo) InitTables(ctx context.Context, tables map[string]string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	for name, query := range tables {
		r.logger.Info("init table started", zap.String("name", name))
		if _, err := r.db.ExecContext(ctx, query); err != nil {
			r.logger.Error("init table failed", zap.String("name", name), zap.Error(err))
			return fmt.Errorf("table %s: %w", name, err)
		}
		r.logger.Info("init table finished", zap.String("name", name))
	}

	return nil
}

func NewPostgresRepo(ctx context.Context, connString string, logger *zap.Logger) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	pr := &PostgresRepo{
		db:     db,
		logger: logger.Named("pgRepo"),
	}

	if err := pr.InitTables(ctx, defaultTables); err != nil {
		return nil, err
	}

	return pr, nil
}
