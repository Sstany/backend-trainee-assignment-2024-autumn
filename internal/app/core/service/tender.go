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

type TenderService struct {
	userRepo         port.UserRepo
	organizationRepo port.OrganizationRepo
	tenderRepo       port.TenderRepo
}

func NewTenderService(repo port.TenderRepo, userRepo port.UserRepo, orgRepo port.OrganizationRepo) *TenderService {
	return &TenderService{
		tenderRepo:       repo,
		userRepo:         userRepo,
		organizationRepo: orgRepo,
	}
}

func (r *TenderService) Create(ctx context.Context, tender *entity.Tender, username string) error {
	tender.ID = entity.TenderID(uuid.NewString())
	tender.Status = entity.TenderStatus(Created)
	tender.CreatedAt = time.Now()

	user, err := r.userRepo.FindUserId(ctx, username)
	if err != nil {
		return err
	}

	if user == "" {
		return ErrUserNotExists
	}

	// TODO check if organization exists

	if err := r.tenderRepo.Create(ctx, tender); err != nil {
		return fmt.Errorf("create tender: %w", err)
	}

	return nil
}

func (r *TenderService) List(ctx context.Context, tenderType []entity.TenderServiceType, limitOffset *entity.RequestLimitOffset) ([]*entity.Tender, error) {
	if len(tenderType) == 0 {
		tenderType = append(tenderType, entity.TenderServiceTypeConstruction)
		tenderType = append(tenderType, entity.TenderServiceTypeDelivery)
		tenderType = append(tenderType, entity.TenderServiceTypeManufacture)
	}

	tenders, err := r.tenderRepo.List(ctx, tenderType, limitOffset)
	if err != nil {
		return nil, fmt.Errorf("list tender: %w", err)
	}

	if len(tenders) == 0 {
		return nil, ErrTenderNotFound
	}

	return tenders, nil
}

func (r *TenderService) ListMy(ctx context.Context, userName string, limitOffset *entity.RequestLimitOffset) ([]*entity.Tender, error) {
	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return nil, ErrUserNotExists
	}

	organization, err := r.organizationRepo.FindOrganizationsByResponsibleUserID(ctx, entity.UserID(userID))
	if err != nil {
		return nil, ErrNotEnoughRights
	}

	tenders, err := r.tenderRepo.ListMy(ctx, organization, limitOffset)
	if err != nil {
		return nil, err
	}

	if len(tenders) == 0 {
		return nil, ErrTenderOrBidNotFound
	}

	return tenders, nil
}

func (r *TenderService) GetUserRights(ctx context.Context, tenderID entity.TenderID, userName string) {

}

func (r *TenderService) GetStatus(ctx context.Context, tenderID entity.TenderID, userName string) (entity.TenderStatus, error) {
	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return "", ErrUserNotExists
	}

	var tender *entity.Tender

	tender, err = r.tenderRepo.Read(ctx, tenderID)
	if err != nil {
		return "", err
	}

	if tender == nil {
		return "", ErrTenderNotFound
	}

	if !r.userRepo.Exists(ctx, userID) {
		return "", ErrUserNotExists
	}

	users, err := r.organizationRepo.FindResponsibleUsers(ctx, []entity.OrganizationID{tender.OrganizationID})
	if err != nil {
		return "", err
	}

	if !slices.Contains(users, entity.UserID(userID)) {
		return "", ErrNotEnoughRights
	}

	return tender.Status, nil
}

func (r *TenderService) SetStatus(ctx context.Context, tenderID entity.TenderID, userName string, status entity.TenderStatus) (*entity.Tender, error) {
	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return nil, ErrUserNotExists
	}

	var tender *entity.Tender

	tender, err = r.tenderRepo.Read(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	if tender == nil {
		return nil, ErrTenderNotFound
	}

	users, err := r.organizationRepo.FindResponsibleUsers(ctx, []entity.OrganizationID{tender.OrganizationID})
	if err != nil {
		return nil, err
	}

	if !slices.Contains(users, entity.UserID(userID)) {
		return nil, ErrNotEnoughRights
	}

	tender.Status = status

	if err := r.tenderRepo.UpdateStatus(ctx, tenderID, status); err != nil {
		return nil, err
	}

	return tender, nil
}

func (r *TenderService) Edit(ctx context.Context, tenderID entity.TenderID, update *entity.TenderUpdate, userName string) (*entity.Tender, error) {
	if entity.NewTenderServiceType(update.ServiceType) == entity.TenderServiceTypeUndefined {
		return nil, ErrWrongInputFormat
	}

	userID, err := r.userRepo.FindUserId(ctx, userName)
	if err != nil {
		return nil, err
	}

	if userID == "" {
		return nil, ErrUserNotExists
	}

	tender, err := r.tenderRepo.Read(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	if tender == nil {
		return nil, ErrTenderNotFound
	}

	orgIDs, err := r.organizationRepo.FindOrganizationsByResponsibleUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(orgIDs) == 0 || !slices.Contains(orgIDs, tender.OrganizationID) {
		return nil, ErrNotEnoughRights
	}

	if err := r.tenderRepo.Update(ctx, tenderID, update); err != nil {
		return nil, err
	}

	return tender.Apply(update), nil
}
