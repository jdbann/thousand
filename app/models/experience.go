package models

import (
	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

// OldExperience holds the details of a single experience.
type OldExperience string

// Experience holds the domain level representation of a vampire's experience.
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
