package repo

import (
	"context"
	"database/sql"
	"errors"

	"avito2024/internal/app/core/entity"

	"go.uber.org/zap"
)

type BidRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

const (
	queryInitBid = `CREATE TABLE IF NOT EXISTS bid (
	id UUID PRIMARY KEY,
	name VARCHAR(100),
	description VARCHAR(500),
	status TEXT,
	tender_id VARCHAR(100),
	author_type TEXT,
	author_id VARCHAR(100),
	version integer DEFAULT 1,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`
	queryCreateBid  = `INSERT INTO bid VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	queryReadMyBids = `SELECT * FROM bid WHERE author_id = $1 ORDER BY name`

	queryReadTenderBids = `SELECT * FROM bid WHERE tender_id = $1`
	queryFindBidAuthor  = `SELECT author_id FROM bid WHERE id = $1`
	queryReadBidByID    = `SELECT * FROM bid WHERE id =$1`

	queryChangeTenderStatus = `UPDATE status FROM tenders WHERE id = $1 SET status = $2`
)

var bidTables map[string]string = map[string]string{
	"bid": queryInitBid,
}

func (r *BidRepo) Create(ctx context.Context, bid *entity.Bid) error {
	_, err := r.db.ExecContext(
		ctx,
		queryCreateBid,
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)

	return err
}

func (r *BidRepo) ReadMyBids(ctx context.Context, userID entity.UserID) ([]*entity.Bid, error) {
	rows, err := r.db.QueryContext(ctx, queryReadMyBids, userID)
	if err != nil {
		return nil, err
	}

	var bids []*entity.Bid

	for rows.Next() {
		bid := new(entity.Bid)

		if err := rows.Scan(
			&bid.ID,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderID,
			&bid.AuthorType,
			&bid.AuthorID,
			&bid.Version,
			&bid.CreatedAt,
		); err != nil {
			return nil, err
		}

		bids = append(bids, bid)
	}
	return bids, nil

}

func (r *BidRepo) ReadBidByID(ctx context.Context, bidID entity.BidId) (*entity.Bid, error) {
	var bid entity.Bid
	if err := r.db.QueryRowContext(ctx, queryReadBidByID, bidID).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &bid, nil
}

func (r *BidRepo) ChangeTenderStatus(ctx context.Context, tenderID entity.TenderID, status entity.TenderStatus) error {
	_, err := r.db.ExecContext(ctx, queryChangeTenderStatus, tenderID, status)
	if err != nil {
		return err
	}
	return nil
}

func (r *BidRepo) ReadTenderBids(ctx context.Context, tenderId entity.TenderID) ([]*entity.Bid, error) {
	rows, err := r.db.QueryContext(ctx, queryReadTenderBids, tenderId)
	if err != nil {
		return nil, err
	}

	var bids []*entity.Bid

	for rows.Next() {
		bid := new(entity.Bid)

		if err := rows.Scan(
			&bid.ID,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderID,
			&bid.AuthorType,
			&bid.AuthorID,
			&bid.Version,
			&bid.CreatedAt,
		); err != nil {
			return nil, err
		}

		bids = append(bids, bid)
	}
	return bids, nil
}

func (r *BidRepo) ReadBidResponsibleUsers(ctx context.Context, bidIDs []entity.BidId) ([]entity.UserID, error) {
	var userIDs []entity.UserID

	for _, bidID := range bidIDs {
		rows, err := r.db.QueryContext(ctx, queryFindBidAuthor, bidID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				r.logger.Error("author not found", zap.Any("bidIds", bidIDs))
				continue
			}

			return nil, err
		}

		for rows.Next() {
			var user entity.UserID
			if err := rows.Scan(&user); err != nil {
				return nil, err
			}
			userIDs = append(userIDs, user)
		}
	}

	if len(userIDs) == 0 {
		return nil, nil
	}

	return userIDs, nil
}

func (r *PostgresRepo) NewBidRepo(ctx context.Context) (*BidRepo, error) {
	ur := &BidRepo{
		db:     r.db,
		logger: r.logger.Named("bid"),
	}

	if err := r.InitTables(ctx, bidTables); err != nil {
		return nil, err
	}

	return ur, nil
}
