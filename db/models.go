// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CharacterType string

const (
	CharacterTypeMortal   CharacterType = "mortal"
	CharacterTypeImmortal CharacterType = "immortal"
)

func (e *CharacterType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CharacterType(s)
	case string:
		*e = CharacterType(s)
	default:
		return fmt.Errorf("unsupported scan type for CharacterType: %T", src)
	}
	return nil
}

type Character struct {
	ID        uuid.UUID
	VampireID uuid.UUID
	Name      string
	Type      CharacterType
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type Experience struct {
	ID          uuid.UUID
	MemoryID    uuid.UUID
	Description string
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}

type Memory struct {
	ID        uuid.UUID
	VampireID uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type Resource struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
	Stationary  bool
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}

type Skill struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}

type Vampire struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
