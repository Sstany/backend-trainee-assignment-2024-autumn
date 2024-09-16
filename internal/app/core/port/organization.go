package port

import (
	"context"

	"avito2024/internal/app/core/entity"
)

type OrganizationRepo interface {
	Exists(context.Context, entity.OrganizationID) bool

	ReadResponsibleUserOrganization(context.Context, entity.UserID) ([]entity.OrganizationID, error)
	FindOrganizationsByResponsibleUserID(context.Context, entity.UserID) ([]entity.OrganizationID, error)
	FindResponsibleUsers(context.Context, []entity.OrganizationID) ([]entity.UserID, error)
}
