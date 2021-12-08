package repository

import (
	"context"
	"errors"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// CreateResource attempts to add a new resource to the DB for the provided
// vampire.
func (m *Repository) CreateResource(ctx context.Context, vampireID uuid.UUID, params models.CreateResourceParams) (models.Resource, error) {
	dbParams := queries.CreateResourceParams{
		VampireID:   vampireID,
		Description: params.Description,
		Stationary:  params.Stationary,
	}

	dbResource, err := m.queries.CreateResource(ctx, dbParams)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.ForeignKeyViolation && pgErr.ConstraintName == "resources_vampire_id_fkey" {
			return models.Resource{}, models.ErrNotFound.Cause(err)
		}

		return models.Resource{}, err
	} else if err != nil {
		return models.Resource{}, err
	}

	return newResource(dbResource), nil
}
