package repo

import (
	"context"
	"database/sql"
	"errors"

	"avito2024/internal/app/core/entity"

	"go.uber.org/zap"
)

const (
	queryInitOrganizationType = `CREATE TYPE organization_type AS ENUM (
		'IE',
		'LLC',
		'JSC'
	);`

	queryInitOrganization = `CREATE TABLE IF NOT EXISTS organization (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(100) NOT NULL,
		description TEXT,
		type organization_type,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	queryInitOrganizationResponsible = `CREATE TABLE IF NOT EXISTS organization_responsible (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
		user_id UUID REFERENCES employee(id) ON DELETE CASCADE
	)`

	queryFindOrganizationsByResponsible = `SELECT organization_id FROM organization_responsible WHERE user_id = $1`
	queryFindResponsibleUsers           = `SELECT user_id FROM organization_responsible WHERE organization_id = $1`

	queryFindOrganizationByID = `SELECT id FROM organization WHERE id = $1`
)

var organizationTables = map[string]string{
	// "organizationType":        queryInitOrganizationType,
	"organization":            queryInitOrganization,
	"organizationResponsible": queryInitOrganizationResponsible,
}

type OrganizationRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func (r *OrganizationRepo) FindOrganizationsByResponsibleUserID(ctx context.Context, userID entity.UserID) ([]entity.OrganizationID, error) {
	rows, err := r.db.QueryContext(ctx, queryFindOrganizationsByResponsible, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	var ids []entity.OrganizationID

	for rows.Next() {
		var id entity.OrganizationID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *OrganizationRepo) ReadResponsibleUserOrganization(ctx context.Context, userID entity.UserID) ([]entity.OrganizationID, error) {
	var organizationIDs []entity.OrganizationID

	rows, err := r.db.QueryContext(ctx, queryFindOrganizationsByResponsible, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var organization entity.OrganizationID
		if err := rows.Scan(&organization); err != nil {
			return nil, err
		}
		organizationIDs = append(organizationIDs, organization)
	}

	return organizationIDs, nil
}

func (r *OrganizationRepo) FindResponsibleUsers(ctx context.Context, organizations []entity.OrganizationID) ([]entity.UserID, error) {
	var users []entity.UserID

	for i := range organizations {
		rows, err := r.db.QueryContext(ctx, queryFindResponsibleUsers, organizations[i])
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var userID entity.UserID

			if err := rows.Scan(&userID); err != nil {
				return nil, err
			}
			users = append(users, userID)
		}
	}

	return users, nil
}

func (r *OrganizationRepo) Exists(ctx context.Context, orgID entity.OrganizationID) bool {
	var id string

	err := r.db.QueryRowContext(ctx, queryFindOrganizationByID, orgID).Scan(&id)
	if err != nil {
		return false
	}

	return id == string(orgID)
}

func (r *PostgresRepo) NewOrganizationRepo(ctx context.Context, isTest bool) (*OrganizationRepo, error) {
	or := &OrganizationRepo{
		db:     r.db,
		logger: r.logger.Named("organization"),
	}

	if isTest {
		if err := r.InitTables(ctx, organizationTables); err != nil {
			return nil, err
		}
	}

	return or, nil
}
