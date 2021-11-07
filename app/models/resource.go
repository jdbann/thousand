package models

import (
	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

type Resource struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
	Stationary  bool
}

func newResource(dbResource db.Resource) Resource {
	return Resource{
		ID:          dbResource.ID,
		VampireID:   dbResource.VampireID,
		Description: dbResource.Description,
		Stationary:  dbResource.Stationary,
	}
}

// OldResource holds the details of a resource possessed by a Vampire.
type OldResource struct {
	ID          int
	Description string `form:"description"`
	Stationary  bool   `form:"stationary"`
	Lost        bool   `form:"lost"`
}
