package models

import (
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
)

type Experience struct {
	ID          uuid.UUID
	MemoryID    uuid.UUID
	Description string
}

func newExperience(experience queries.Experience) Experience {
	return Experience{
		ID:          experience.ID,
		MemoryID:    experience.MemoryID,
		Description: experience.Description,
	}
}
