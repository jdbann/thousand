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

// CreateSkill attempts to add a new skill to the DB for the provided vampire.
func (m *Repository) CreateSkill(ctx context.Context, vampireID uuid.UUID, description string) (models.Skill, error) {
	params := queries.CreateSkillParams{
		VampireID:   vampireID,
		Description: description,
	}

	dbSkill, err := m.queries.CreateSkill(ctx, params)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.ForeignKeyViolation && pgErr.ConstraintName == "skills_vampire_id_fkey" {
			return models.Skill{}, models.ErrNotFound.Cause(err)
		}

		return models.Skill{}, err
	} else if err != nil {
		return models.Skill{}, err
	}

	return newSkill(dbSkill), nil
}
