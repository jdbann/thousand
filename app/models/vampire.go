package models

import "errors"

// ErrNotFound is returned when trying to reference a model by an ID if a model
// cannot be found with that ID.
var ErrNotFound = errors.New("Record not found")

// Vampire holds everything related to a vampire.
type Vampire struct {
	Details    *Details
	Memories   [5]Memory
	Skills     []Skill
	Resources  []Resource
	Characters []Character
}

// AddExperience adds an experience to the indicated memory.
func (v *Vampire) AddExperience(memoryID int, experienceString string) error {
	experience := Experience(experienceString)

	if memoryID < 0 || memoryID >= len(v.Memories) {
		return ErrNotFound
	}

	return v.Memories[memoryID].AddExperience(experience)
}

// AddSkill adds an unchecked skill to the Vampire.
func (v *Vampire) AddSkill(skill *Skill) {
	skill.ID = len(v.Skills) + 1
	v.Skills = append(v.Skills, *skill)
}

// FindSkill retrieves a skill based on an ID from the Vampire's list of
// skills.
func (v *Vampire) FindSkill(skillID int) (*Skill, error) {
	for _, skill := range v.Skills {
		if skill.ID == skillID {
			return &skill, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateSkill replaces a Vampire's existing skill with the new one based on
// the new skill's ID.
func (v *Vampire) UpdateSkill(newSkill *Skill) error {
	var skills []Skill
	found := false

	for _, originalSkill := range v.Skills {
		if originalSkill.ID == newSkill.ID {
			found = true
			skills = append(skills, *newSkill)
		} else {
			skills = append(skills, originalSkill)
		}
	}

	if found {
		v.Skills = skills
		return nil
	}

	return ErrNotFound
}

// AddResource adds a resource to the Vampire.
func (v *Vampire) AddResource(resource *Resource) {
	resource.ID = len(v.Resources) + 1
	v.Resources = append(v.Resources, *resource)
}

// FindResource retrieves a resource based on an ID from the Vampire's list of
// resources.
func (v *Vampire) FindResource(resourceID int) (*Resource, error) {
	for _, resource := range v.Resources {
		if resource.ID == resourceID {
			return &resource, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateResource replaces a Vampire's existing resource with the new one
// based on the new resource's ID.
func (v *Vampire) UpdateResource(newResource *Resource) error {
	var resources []Resource
	found := false

	for _, originalResource := range v.Resources {
		if originalResource.ID == newResource.ID {
			found = true
			resources = append(resources, *newResource)
		} else {
			resources = append(resources, originalResource)
		}
	}

	if found {
		v.Resources = resources
		return nil
	}

	return ErrNotFound
}

// AddCharacter adds a character to the Vampire.
func (v *Vampire) AddCharacter(character *Character) {
	character.ID = len(v.Characters) + 1
	v.Characters = append(v.Characters, *character)
}
