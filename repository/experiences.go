package repository

import (
	"context"
	"errors"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// CreateExperience attempts to add a new experience to the DB for the provided
// memory.
func (m *Repository) CreateExperience(ctx context.Context, vampireID, memoryID uuid.UUID, description string) (models.Experience, error) {
	params := queries.CreateExperienceParams{
		VampireID:   vampireID,
		MemoryID:    memoryID,
		Description: description,
	}

	dbExperience, err := m.queries.CreateExperience(ctx, params)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.NotNullViolation && pgErr.ColumnName == "memory_id" {
			return models.Experience{}, models.ErrNotFound.Cause(err)
		}

		if pgErr.Code == models.PgErrCodeMemoryFull {
			return models.Experience{}, models.ErrMemoryFull.Cause(err)
		}

		return models.Experience{}, err
	} else if err != nil {
		return models.Experience{}, err
	}

	return newExperience(dbExperience), nil
}

// GetExperiences attempts to retrieve all the experiences from the DB for the
// provided vampire.
func (m *Repository) GetExperiences(ctx context.Context, vampireID uuid.UUID) ([]models.Experience, error) {
	dbExperiences, err := m.queries.GetExperiencesForVampire(ctx, vampireID)
	if err != nil {
		return nil, err
	}

	var experiences []models.Experience
	for _, dbExperience := range dbExperiences {
		experiences = append(experiences, newExperience(dbExperience))
	}

	return experiences, nil
}
