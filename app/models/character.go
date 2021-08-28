package models

import "errors"

// ErrNotFound is returned when trying to reference a model by an ID if a model
// cannot be found with that ID.
var ErrNotFound = errors.New("Record not found")

// Character holds everything related to a character.
type Character struct {
	Details  *Details
	Memories [5]Memory
}

// AddExperience adds an experience to the indicated memory.
func (c *Character) AddExperience(memoryID int, experienceString string) error {
	experience := Experience(experienceString)

	if memoryID < 0 || memoryID >= len(c.Memories) {
		return ErrNotFound
	}

	return c.Memories[memoryID].AddExperience(experience)
}
