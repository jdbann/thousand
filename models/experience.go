package models

import (
	"github.com/google/uuid"
)

type Experience struct {
	ID          uuid.UUID
	MemoryID    uuid.UUID
	Description string
}
