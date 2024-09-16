package repo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"avito2024/internal/app/core/entity"
)

type TenderRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

const (
	queryInitTender = `CREATE TABLE IF NOT EXISTS tenders (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(100) NOT NULL,
		description VARCHAR(500),
		service_type TEXT,
		status TEXT,
		organization_id VARCHAR(50) UNIQUE NOT NULL,
		version INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	queryCreateTender  = `INSERT INTO tenders VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	queryListTenders   = `SELECT * FROM tenders WHERE service_type = ANY($1) AND status = 'Published' ORDER BY name `
	queryListMyTenders = `SELECT * FROM tenders WHERE organization_id = ANY($1) ORDER BY name `

	orderByName             = ` ORDER BY name`
	queryReadTender         = `SELECT * FROM tenders WHERE id = $1`
	queryUpdateTenderStatus = `UPDATE tenders SET status = $1 WHERE id = $2`
)

var tenderTables map[string]string = map[string]string{
	"tender": queryInitTender,
}

func (r *TenderRepo) Create(ctx context.Context, tender *entity.Tender) error {
	_, err := r.db.ExecContext(
		ctx,
		queryCreateTender,
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationID,
		&tender.Version,
		&tender.CreatedAt,
	)

	return err
}

func (r *TenderRepo) Read(ctx context.Context, tenderID entity.TenderID) (*entity.Tender, error) {
	var tender entity.Tender

	err := r.db.QueryRowContext(ctx, queryReadTender, tenderID).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationID,
		&tender.Version,
		&tender.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("no tender found", zap.String("id", string(tenderID)))
			return nil, nil
		}

		return nil, err
	}

	return &tender, nil
}

func (r *TenderRepo) List(ctx context.Context, tenderTypes []entity.TenderServiceType, limitOffset *entity.RequestLimitOffset) ([]*entity.Tender, error) {
	var tenders []*entity.Tender

	query, args := buildLimitOffset(
		queryListTenders,
		[]any{
			pq.Array(tenderTypes),
		},
		limitOffset,
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	for rows.Next() {
		var tender entity.Tender

		if err := rows.Scan(
			&tender.ID,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationID,
			&tender.Version,
			&tender.CreatedAt,
		); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (r *TenderRepo) ListMy(ctx context.Context, organizations []entity.OrganizationID, limitOffset *entity.RequestLimitOffset) ([]*entity.Tender, error) {
	query, args := buildLimitOffset(
		queryListMyTenders,
		[]any{
			pq.Array(organizations),
		},
		limitOffset,
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	var tenders []*entity.Tender

	for rows.Next() {
		var tender entity.Tender

		if err := rows.Scan(
			&tender.ID,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationID,
			&tender.Version,
			&tender.CreatedAt,
		); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (r *TenderRepo) UpdateStatus(ctx context.Context, tenderID entity.TenderID, tenderStatus entity.TenderStatus) error {
	_, err := r.db.ExecContext(ctx, queryUpdateTenderStatus, tenderStatus, tenderID)

	return err
}

func (r *TenderRepo) Update(ctx context.Context, tenderID entity.TenderID, update *entity.TenderUpdate) error {
	query := sqlbuilder.Update("tenders")

	if update == nil {
		return nil
	}

	var assign []string

	if update.Name != "" {
		assign = append(assign, query.Assign("name", update.Name))
	}

	if update.Description != "" {
		assign = append(assign, query.Assign("description", update.Description))
	}

	if update.ServiceType != "" {
		assign = append(assign, query.Assign("service_type", update.ServiceType))
	}

	query.Set(assign...).Where(query.Equal("id", tenderID))

	query.SetFlavor(sqlbuilder.PostgreSQL)

	queryString, args := query.Build()
	stmt, err := r.db.PrepareContext(ctx, queryString)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, args...)

	return err
}

func (r *PostgresRepo) NewTenderRepo(ctx context.Context) (*TenderRepo, error) {
	tr := &TenderRepo{
		db:     r.db,
		logger: r.logger.Named("tender"),
	}

	if err := r.InitTables(ctx, tenderTables); err != nil {
		return nil, err
	}

	return tr, nil
}

func buildLimitOffset(query string, args []any, limitOffset *entity.RequestLimitOffset) (string, []any) {
	var s strings.Builder
	s.WriteString(query)

	var offsetSet bool

	if limitOffset.Offset > 0 {
		args = append(args, limitOffset.Offset)
		s.WriteString("OFFSET $2")
		s.WriteString(" ")
		offsetSet = true
	}

	if limitOffset.Limit > 0 {
		args = append(args, limitOffset.Limit)
		s.WriteString("LIMIT $")
		if offsetSet {
			s.WriteString("3")
		} else {
			s.WriteString("2")
		}
	}

	return s.String(), args
}
