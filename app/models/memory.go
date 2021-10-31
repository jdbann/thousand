package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Memory holds a maximum of three experiences.
type Memory struct {
	ID          int
	Experiences []Experience
}

// NewMemory will replace Memory when the DB persistence work is complete.
// TODO: Replace Memory with NewMemory
type NewMemory struct {
	ID        uuid.UUID
	VampireID uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type WholeMemory struct {
	NewMemory
	Experiences []NewExperience
}

// Full returns true if there is no more room for experiences in this memory.
func (m *Memory) Full() bool {
	return len(m.Experiences) >= 3
}

// AddExperience adds a new experience into a Memory if there is at least one
// space available.
func (m *Memory) AddExperience(experience Experience) error {
	if len(m.Experiences) >= 3 {
		return ErrMemoryFull
	}

	m.Experiences = append(m.Experiences, Experience(experience))
	return nil
}
