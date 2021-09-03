package models

import "errors"

// ErrNotFound is returned when trying to reference a model by an ID if a model
// cannot be found with that ID.
var ErrNotFound = errors.New("Record not found")

// Vampire holds everything related to a vampire.
type Vampire struct {
	Details    *Details
	Memories   []*Memory
	Skills     []*Skill
	Resources  []*Resource
	Characters []*Character
	Marks      []Mark
}

func (v *Vampire) findMemory(memoryID int) (*Memory, error) {
	for _, memory := range v.Memories {
		if memory.ID == memoryID {
			return memory, nil
		}
	}

	return nil, ErrNotFound
}

// AddExperience adds an experience to the indicated memory.
func (v *Vampire) AddExperience(memoryID int, experienceString string) error {
	memory, err := v.findMemory(memoryID)
	if err != nil {
		return err
	}

	experience := Experience(experienceString)
	return memory.AddExperience(experience)
}

// AddSkill adds an unchecked skill to the Vampire.
func (v *Vampire) AddSkill(skill *Skill) {
	skill.ID = len(v.Skills) + 1
	v.Skills = append(v.Skills, skill)
}

// FindSkill retrieves a skill based on an ID from the Vampire's list of
// skills.
func (v *Vampire) FindSkill(skillID int) (*Skill, error) {
	for _, skill := range v.Skills {
		if skill.ID == skillID {
			return skill, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateSkill replaces a Vampire's existing skill with the new one based on
// the new skill's ID.
func (v *Vampire) UpdateSkill(newSkill *Skill) error {
	oldSkill, err := v.FindSkill(newSkill.ID)
	if err != nil {
		return err
	}

	*oldSkill = *newSkill

	return nil
}

// AddResource adds a resource to the Vampire.
func (v *Vampire) AddResource(resource *Resource) {
	resource.ID = len(v.Resources) + 1
	v.Resources = append(v.Resources, resource)
}

// FindResource retrieves a resource based on an ID from the Vampire's list of
// resources.
func (v *Vampire) FindResource(resourceID int) (*Resource, error) {
	for _, resource := range v.Resources {
		if resource.ID == resourceID {
			return resource, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateResource replaces a Vampire's existing resource with the new one
// based on the new resource's ID.
func (v *Vampire) UpdateResource(newResource *Resource) error {
	oldResource, err := v.FindResource(newResource.ID)
	if err != nil {
		return err
	}

	*oldResource = *newResource

	return nil
}

// AddCharacter adds a character to the Vampire.
func (v *Vampire) AddCharacter(character *Character) {
	character.ID = len(v.Characters) + 1
	v.Characters = append(v.Characters, character)
}

// FindCharacter retrieves a character based on an ID from the Vampire's list of
// characters.
func (v *Vampire) FindCharacter(characterID int) (*Character, error) {
	for _, character := range v.Characters {
		if character.ID == characterID {
			return character, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateCharacter replaces a Vampire's existing character with the new one
// based on the new character's ID.
func (v *Vampire) UpdateCharacter(newCharacter *Character) error {
	oldCharacter, err := v.FindCharacter(newCharacter.ID)
	if err != nil {
		return err
	}

	*oldCharacter = *newCharacter

	return nil
}

// AddDescriptor adds a descriptor to the indicated character.
func (v *Vampire) AddDescriptor(characterID int, descriptor string) error {
	character, err := v.FindCharacter(characterID)
	if err != nil {
		return err
	}

	character.AddDescriptor(descriptor)

	return nil
}

// AddMark adds a mark to the Vampire.
func (v *Vampire) AddMark(mark *Mark) {
	mark.ID = len(v.Marks) + 1
	v.Marks = append(v.Marks, *mark)
}
