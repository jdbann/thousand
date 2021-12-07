package models

import (
	"github.com/google/uuid"
)

type Character struct {
	ID   uuid.UUID
	Name string
	Type string
}

type CreateCharacterParams struct {
	Name string `form:"name"`
	Type string `form:"type"`
}
