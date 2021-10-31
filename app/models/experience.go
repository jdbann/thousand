package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Experience holds the details of a single experience.
type Experience string

// NewExperience will replace Experience when the DB persistence work is
// complete.
// TODO: Replace Experience with NewExperience
type NewExperience struct {
	ID          uuid.UUID
	MemoryID    uuid.UUID
	Description string
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}
