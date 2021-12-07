package models

import (
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
)

type Memory struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Experiences []Experience
}

func newMemory(dbMemory queries.Memory, dbExperiences []queries.Experience) Memory {
	var experiences = make([]Experience, len(dbExperiences))
	for i, dbExperience := range dbExperiences {
		experiences[i] = newExperience(dbExperience)
	}

	return Memory{
		ID:          dbMemory.ID,
		VampireID:   dbMemory.VampireID,
		Experiences: experiences,
	}
}

// Full returns true if there is no more room for additional experiences in this
// memory.
func (m Memory) Full() bool {
	return len(m.Experiences) >= 3
}
