// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Vampire struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
