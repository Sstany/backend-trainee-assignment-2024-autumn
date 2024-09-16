package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"avito2024/internal/app/core/entity"
	"avito2024/internal/app/core/port"
)

type BidService struct {
	userRepo         port.UserRepo
	organizationRepo port.OrganizationRepo
	bidRepo          port.BidRepo
	tenderRepo       port.TenderRepo
}

func NewBidService(
	bidRepo port.BidRepo,
	userRepo port.UserRepo,
	organizationRepo port.OrganizationRepo,
	tenderRepo port.TenderRepo,
) *BidService {
	return &BidService{
		bidRepo:          bidRepo,
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
		tenderRepo:       tenderRepo,
	}
}

func (r *BidService) Create(ctx context.Context, bid *entity.Bid) error {
	bid.ID = entity.BidId(uuid.NewString())
	bid.Status = entity.BidStatus(Created)
	bid.CreatedAt = time.Now()

	switch bid.AuthorType {
	case entity.BidAuthorOrganization:
		exists := r.organizationRepo.Exists(ctx, entity.OrganizationID(bid.AuthorID))
		if !exists {
			return ErrUserNotExists
		}
	case entity.BidAuthorUser:
		exists := r.userRepo.Exists(ctx, entity.UserID(bid.AuthorID))
		if !exists {
			return ErrUserNotExists
		}
	default:
		return ErrWrongInputFormat
	}

	tender, err := r.tenderRepo.Read(ctx, bid.TenderID)
	if err != nil {
		return err
	}

	if tender == nil {
		return ErrTenderNotFound
	}

	if err := r.bidRepo.Create(ctx, bid); err != nil {
		return fmt.Errorf("create bid: %w", err)
	}

	return nil
}

func (r *BidService) ListBidsMy(ctx context.Context, userName string) ([]*entity.Bid, error) {
	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return nil, ErrUserNotExists
	}

	bids, err := r.bidRepo.ReadMyBids(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list bids: %w", err)
	}

	return bids, nil
}

func (r *BidService) ListTenderBids(ctx context.Context, tenderID entity.TenderID, userName string) ([]*entity.Bid, error) {
	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return nil, ErrUserNotExists
	}

	tender, err := r.tenderRepo.Read(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	if tender == nil {
		return nil, ErrTenderOrBidNotFound
	}

	orgIDs, err := r.organizationRepo.ReadResponsibleUserOrganization(ctx, userID)
	if err != nil {
		return nil, err
	}

	tenderOrgID := tender.OrganizationID

	if !slices.Contains(orgIDs, tenderOrgID) {
		return nil, ErrNotEnoughRights
	}

	bids, err := r.bidRepo.ReadTenderBids(ctx, tenderID)
	if err != nil {
		return nil, ErrTenderOrBidNotFound
	}

	return bids, nil
}

func (r *BidService) SubmitDecision(ctx context.Context, bidID entity.BidId, decision string, userName string) (*entity.Bid, error) {
	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return nil, ErrUserNotExists
	}

	bid, err := r.bidRepo.ReadBidByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	if bid == nil {
		return nil, ErrBidNotFound
	}

	users, err := r.bidRepo.ReadBidResponsibleUsers(ctx, []entity.BidId{bidID})
	if err != nil {
		return nil, ErrBidNotFound
	}

	if !slices.Contains(users, entity.UserID(userID)) {
		return nil, ErrNotEnoughRights
	}

	tender, err := r.tenderRepo.Read(ctx, bid.TenderID)
	if err != nil {
		return nil, err
	}

	if tender == nil {
		return nil, ErrTenderNotFound
	}

	switch entity.BidReviewDecision(decision) {
	case entity.BidReviewApproved:
		tender.Status = entity.TenderStatusClosed
		r.tenderRepo.UpdateStatus(ctx, tender.ID, tender.Status)
	}

	return bid, nil
}

func (r *BidService) Status(ctx context.Context, bidID entity.BidId, userName string) (entity.BidStatus, error) {
	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return "", ErrUserNotExists
	}
	var bids []entity.BidId
	bids = append(bids, bidID)

	users, err := r.bidRepo.ReadBidResponsibleUsers(ctx, bids)
	if err != nil {
		return "", ErrBidNotFound
	}

	if !slices.Contains(users, entity.UserID(userID)) {
		return "", ErrNotEnoughRights
	}

	bid, err := r.bidRepo.ReadBidByID(ctx, bidID)
	if err != nil {
		return "", ErrBidNotFound
	}

	return bid.Status, nil

}
