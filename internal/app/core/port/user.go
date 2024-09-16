package port

import (
	"context"

	"avito2024/internal/app/core/entity"
)

type UserRepo interface {
	FindUserId(ctx context.Context, userName string) (entity.UserID, error)

	Exists(context.Context, entity.UserID) bool
}
