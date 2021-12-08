package models

import (
	"github.com/google/uuid"
)

type Mark struct {
	ID          uuid.UUID
	Description string
}
