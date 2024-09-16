package port

import (
	"context"

	"avito2024/internal/app/core/entity"
)

type TenderRepo interface {
	Create(context.Context, *entity.Tender) error
	List(context.Context, []entity.TenderServiceType, *entity.RequestLimitOffset) ([]*entity.Tender, error)
	Read(context.Context, entity.TenderID) (*entity.Tender, error)
	UpdateStatus(context.Context, entity.TenderID, entity.TenderStatus) error
	ListMy(context.Context, []entity.OrganizationID, *entity.RequestLimitOffset) ([]*entity.Tender, error)
	Update(context.Context, entity.TenderID, *entity.TenderUpdate) error
}
