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

func (m *Repository) CreateMark(ctx context.Context, vampireID uuid.UUID, description string) (models.Mark, error) {
	params := queries.CreateMarkParams{
		VampireID:   vampireID,
		Description: description,
	}

	dbMark, err := m.queries.CreateMark(ctx, params)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.ForeignKeyViolation && pgErr.ConstraintName == "marks_vampire_id_fkey" {
			return models.Mark{}, models.ErrNotFound.Cause(err)
		}

		return models.Mark{}, err
	} else if err != nil {
		return models.Mark{}, err
	}

	return newMark(dbMark), nil
}
