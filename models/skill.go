package models

import (
	"github.com/google/uuid"
)

type Skill struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
}
