package models

import (
	"github.com/google/uuid"
)

type Memory struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Experiences []Experience
}

// Full returns true if there is no more room for additional experiences in this
// memory.
func (m Memory) Full() bool {
	return len(m.Experiences) >= 3
}
