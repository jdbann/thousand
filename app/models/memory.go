package models

import (
	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

// OldMemory holds a maximum of three experiences.
type OldMemory struct {
	ID          int
	Experiences []Experience
}

// Memory holds the domain level representation of a vampire's memory.
type Memory struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Experiences []NewExperience
}

func newMemory(dbMemory db.Memory, dbExperiences []db.Experience) Memory {
	var experiences = make([]NewExperience, len(dbExperiences))
	for i, experience := range dbExperiences {
		experiences[i] = NewExperience(experience)
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
func (m *OldMemory) AddExperience(experience Experience) error {
	if len(m.Experiences) >= 3 {
		return ErrMemoryFull
	}

	m.Experiences = append(m.Experiences, Experience(experience))
	return nil
}
