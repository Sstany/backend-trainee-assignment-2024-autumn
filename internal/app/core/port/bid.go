package port

import (
	"context"

	"avito2024/internal/app/core/entity"
)

type BidRepo interface {
	Create(context.Context, *entity.Bid) error
	ReadMyBids(context.Context, entity.UserID) ([]*entity.Bid, error)
	ReadTenderBids(context.Context, entity.TenderID) ([]*entity.Bid, error)
	ReadBidResponsibleUsers(context.Context, []entity.BidId) ([]entity.UserID, error)
	ReadBidByID(context.Context, entity.BidId) (*entity.Bid, error)
}
