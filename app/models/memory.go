package models

import (
	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

// OldMemory holds a maximum of three experiences.
type OldMemory struct {
	ID          int
	Experiences []OldExperience
}

// Memory holds the domain level representation of a vampire's memory.
type Memory struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Experiences []Experience
}

func newMemory(dbMemory db.Memory, dbExperiences []db.Experience) Memory {
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

// Full returns true if there is no more room for experiences in this memory.
func (m *OldMemory) Full() bool {
	return len(m.Experiences) >= 3
}

// AddExperience adds a new experience into a Memory if there is at least one
// space available.
func (m *OldMemory) AddExperience(experience OldExperience) error {
	if len(m.Experiences) >= 3 {
		return ErrMemoryFull
	}

	m.Experiences = append(m.Experiences, OldExperience(experience))
	return nil
}
