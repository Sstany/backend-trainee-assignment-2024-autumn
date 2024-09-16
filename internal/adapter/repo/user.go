package repo

import (
	"context"
	"database/sql"
	"errors"

	"avito2024/internal/app/core/entity"

	"go.uber.org/zap"
)

const (
	queryInitEmploee = `CREATE TABLE IF NOT EXISTS employee (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		username VARCHAR(50) UNIQUE NOT NULL,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	queryUserIDByUsername = "SELECT id FROM employee WHERE username = $1"
	queryUserIDByID       = "SELECT id FROM employee WHERE id = $1"
)

var userTables = map[string]string{
	"employee": queryInitEmploee,
}

type UserRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func (r *UserRepo) FindUserId(ctx context.Context, userName string) (entity.UserID, error) {
	var userID entity.UserID

	if err := r.db.QueryRowContext(ctx, queryUserIDByUsername, userName).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return userID, nil
}

func (r *UserRepo) Exists(ctx context.Context, userID entity.UserID) bool {
	var id string

	err := r.db.QueryRowContext(ctx, queryUserIDByID, userID).Scan(&id)
	if err != nil {
		return false
	}

	return id == string(userID)
}

func (r *PostgresRepo) NewUserRepo(ctx context.Context, isTest bool) (*UserRepo, error) {
	ur := &UserRepo{
		db:     r.db,
		logger: r.logger.Named("tender"),
	}

	if isTest {
		if err := r.InitTables(ctx, userTables); err != nil {
			return nil, err
		}
	}

	return ur, nil
}
