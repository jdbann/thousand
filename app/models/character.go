package models

import "errors"

// ErrNotFound is returned when trying to reference a model by an ID if a model
// cannot be found with that ID.
var ErrNotFound = errors.New("Record not found")

// Character holds everything related to a character.
type Character struct {
	Details   *Details
	Memories  [5]Memory
	Skills    []Skill
	Resources []Resource
}

// AddExperience adds an experience to the indicated memory.
func (c *Character) AddExperience(memoryID int, experienceString string) error {
	experience := Experience(experienceString)

	if memoryID < 0 || memoryID >= len(c.Memories) {
		return ErrNotFound
	}

	return c.Memories[memoryID].AddExperience(experience)
}

// AddSkill adds an unchecked skill to the Character.
func (c *Character) AddSkill(skill *Skill) {
	skill.ID = len(c.Skills) + 1
	c.Skills = append(c.Skills, *skill)
}

// FindSkill retrieves a skill based on an ID from the Character's list of
// skills.
func (c *Character) FindSkill(skillID int) (*Skill, error) {
	for _, skill := range c.Skills {
		if skill.ID == skillID {
			return &skill, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateSkill replaces a Character's existing skill with the new one based on
// the new skill's ID.
func (c *Character) UpdateSkill(newSkill *Skill) error {
	var skills []Skill
	found := false

	for _, originalSkill := range c.Skills {
		if originalSkill.ID == newSkill.ID {
			found = true
			skills = append(skills, *newSkill)
		} else {
			skills = append(skills, originalSkill)
		}
	}

	if found {
		c.Skills = skills
		return nil
	}

	return ErrNotFound
}

// AddResource adds a resource to the Character.
func (c *Character) AddResource(resource *Resource) {
	resource.ID = len(c.Resources) + 1
	c.Resources = append(c.Resources, *resource)
}
