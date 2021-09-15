package models

import (
	"database/sql"
	"errors"
	"time"

	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

// ErrMemoryFull is returned when trying to add experiences to a full memory.
var ErrMemoryFull = errors.New("Memory is full")

// Memory holds a maximum of three experiences.
type Memory struct {
	ID          int
	Experiences []Experience
}

// Ensure that modals.NewMemory is interchangeable with db.Memory at compile
// time with this assignment.
var _ db.Memory = db.Memory(NewMemory{})

// NewMemory holds everything related to a vampire's memory. It will replace
// Memory when the application is modified to store everything on a database
// instead of in memory.
type NewMemory struct {
	ID        uuid.UUID
	VampireID uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
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
