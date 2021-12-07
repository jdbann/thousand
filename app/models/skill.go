package models

import (
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
)

type Skill struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
}

func newSkill(dbSkill queries.Skill) Skill {
	return Skill{
		ID:          dbSkill.ID,
		VampireID:   dbSkill.VampireID,
		Description: dbSkill.Description,
	}
}
