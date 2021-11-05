package models

import (
	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

type Skill struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
}

func newSkill(dbSkill db.Skill) Skill {
	return Skill{
		ID:          dbSkill.ID,
		VampireID:   dbSkill.VampireID,
		Description: dbSkill.Description,
	}
}

// OldSkill holds the details of an ability possessed by a Vampire.
type OldSkill struct {
	ID          int
	Description string `form:"description"`
	Checked     bool   `form:"checked"`
}
