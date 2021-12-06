package models

import (
	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

type Mark struct {
	ID          uuid.UUID
	Description string
}

func newMark(dbMark db.Mark) Mark {
	return Mark{
		ID:          dbMark.ID,
		Description: dbMark.Description,
	}
}
