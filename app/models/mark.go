package models

import (
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
)

type Mark struct {
	ID          uuid.UUID
	Description string
}

func newMark(dbMark queries.Mark) Mark {
	return Mark{
		ID:          dbMark.ID,
		Description: dbMark.Description,
	}
}
