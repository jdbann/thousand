package repository

import (
	"context"
	"errors"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (m *Repository) GetMemory(ctx context.Context, vampireID, id uuid.UUID) (models.Memory, error) {
	params := queries.GetMemoryParams{
		VampireID: vampireID,
		MemoryID:  id,
	}

	dbMemory, err := m.queries.GetMemory(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Memory{}, models.ErrNotFound.Cause(err)
	} else if err != nil {
		return models.Memory{}, err
	}

	// TODO: Also return experiences?
	return newMemory(dbMemory, []queries.Experience{}), nil
}
