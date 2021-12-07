package models

import (
	"github.com/google/uuid"
)

type Resource struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
	Stationary  bool
}

type CreateResourceParams struct {
	Description string `form:"description"`
	Stationary  bool   `form:"stationary"`
}
