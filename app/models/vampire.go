package models

import (
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
)

const (
	// VampireMemorySize specifies how many active memories a vampire should have,
	// whether they are empty or not.
	VampireMemorySize = 5
)

type Vampire struct {
	ID         uuid.UUID
	Name       string
	Memories   []Memory
	Skills     []Skill
	Resources  []Resource
	Characters []Character
	Marks      []Mark
}

func newVampire(dbVampire queries.Vampire, memories []Memory, skills []Skill, resources []Resource, characters []Character, marks []Mark) Vampire {
	return Vampire{
		ID:         dbVampire.ID,
		Name:       dbVampire.Name,
		Memories:   memories,
		Skills:     skills,
		Resources:  resources,
		Characters: characters,
		Marks:      marks,
	}
}
