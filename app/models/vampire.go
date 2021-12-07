package models

import (
	"github.com/google/uuid"
)

// VampireMemorySize specifies how many active memories a vampire should have,
// whether they are empty or not.
const VampireMemorySize = 5

type Vampire struct {
	ID         uuid.UUID
	Name       string
	Memories   []Memory
	Skills     []Skill
	Resources  []Resource
	Characters []Character
	Marks      []Mark
}
