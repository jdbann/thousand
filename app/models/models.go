package models

import (
	"context"

	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

// Models holds a configured *db.Queries which can be used to load values from
// the DB and return them as model types from this package.
type Models struct {
	*db.Queries
}

// DBTX is an alias of db.DBTX to ensure that the models package completely
// encapsulates any DB interactions.
type DBTX = db.DBTX

// NewModels returns a newly configured models.Models struct.
func NewModels(dbtx DBTX) *Models {
	return &Models{db.New(dbtx)}
}

// CreateVampire attempts to create a new vampire in the DB with the provided
// name.
func (m *Models) CreateVampire(ctx context.Context, name string) (NewVampire, error) {
	v, err := m.Queries.CreateVampire(ctx, name)
	if err != nil {
		return NewVampire{}, err
	}

	return NewVampire(v), nil
}

// GetVampire attempts to retrieve a vampire from the DB with the provided ID.
func (m *Models) GetVampire(ctx context.Context, id uuid.UUID) (NewVampire, error) {
	v, err := m.Queries.GetVampire(ctx, id)
	if err != nil {
		return NewVampire{}, err
	}

	return NewVampire(v), nil
}

// GetVampires attempts to retrieve all the vampires from the DB.
func (m *Models) GetVampires(ctx context.Context) ([]NewVampire, error) {
	vs, err := m.Queries.GetVampires(ctx)
	if err != nil {
		return []NewVampire{}, err
	}

	nvs := make([]NewVampire, len(vs))
	for i, v := range vs {
		nvs[i] = NewVampire(v)
	}

	return nvs, nil
}
