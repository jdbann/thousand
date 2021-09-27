// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Memory struct {
	ID        uuid.UUID
	VampireID uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type Vampire struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
