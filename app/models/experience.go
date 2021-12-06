package models

import (
	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

type Experience struct {
	ID          uuid.UUID
	MemoryID    uuid.UUID
	Description string
}

func newExperience(experience db.Experience) Experience {
	return Experience{
		ID:          experience.ID,
		MemoryID:    experience.MemoryID,
		Description: experience.Description,
	}
}
