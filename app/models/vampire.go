package models

import (
	"errors"
	"sort"

	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

// ErrNotFound is returned when trying to reference a model by an ID if a model
// cannot be found with that ID.
var ErrNotFound = errors.New("Record not found")

const (
	// vampireMemorySize specifies how many active memories a vampire should have,
	// whether they are empty or not.
	vampireMemorySize = 5
)

// OldVampire holds everything related to a vampire.
type OldVampire struct {
	Details    *Details
	Memories   []*OldMemory
	Skills     []*OldSkill
	Resources  []*Resource
	Characters []*Character
	Marks      []*Mark
}

// Vampire holds the domain level representation of a vampire.
type Vampire struct {
	ID       uuid.UUID
	Name     string
	Memories []Memory
	Skills   []Skill
}

func newVampire(dbVampire db.Vampire, memories []Memory, skills []Skill) Vampire {
	return Vampire{
		ID:       dbVampire.ID,
		Name:     dbVampire.Name,
		Memories: memories,
		Skills:   skills,
	}
}

func (v *OldVampire) findMemory(memoryID int) (*OldMemory, error) {
	for _, memory := range v.Memories {
		if memory.ID == memoryID {
			return memory, nil
		}
	}

	return nil, ErrNotFound
}

// ForgetMemory replaces an existing memory with a blank memory.
func (v *OldVampire) ForgetMemory(memoryID int) error {
	memory, err := v.findMemory(memoryID)
	if err != nil {
		return err
	}

	memories := make([]*OldMemory, len(v.Memories))
	copy(memories, v.Memories)
	sort.Slice(memories, func(i, j int) bool { return memories[i].ID > memories[j].ID })

	newID := memories[0].ID + 1

	newMemory := OldMemory{
		ID:          newID,
		Experiences: []OldExperience{},
	}

	*memory = newMemory

	return nil
}

// AddExperience adds an experience to the indicated memory.
func (v *OldVampire) AddExperience(memoryID int, experienceString string) error {
	memory, err := v.findMemory(memoryID)
	if err != nil {
		return err
	}

	experience := OldExperience(experienceString)
	return memory.AddExperience(experience)
}

// AddSkill adds an unchecked skill to the Vampire.
func (v *OldVampire) AddSkill(skill *OldSkill) {
	skill.ID = len(v.Skills) + 1
	v.Skills = append(v.Skills, skill)
}

// FindSkill retrieves a skill based on an ID from the Vampire's list of
// skills.
func (v *OldVampire) FindSkill(skillID int) (*OldSkill, error) {
	for _, skill := range v.Skills {
		if skill.ID == skillID {
			return skill, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateSkill replaces a Vampire's existing skill with the new one based on
// the new skill's ID.
func (v *OldVampire) UpdateSkill(newSkill *OldSkill) error {
	oldSkill, err := v.FindSkill(newSkill.ID)
	if err != nil {
		return err
	}

	*oldSkill = *newSkill

	return nil
}

// AddResource adds a resource to the Vampire.
func (v *OldVampire) AddResource(resource *Resource) {
	resource.ID = len(v.Resources) + 1
	v.Resources = append(v.Resources, resource)
}

// FindResource retrieves a resource based on an ID from the Vampire's list of
// resources.
func (v *OldVampire) FindResource(resourceID int) (*Resource, error) {
	for _, resource := range v.Resources {
		if resource.ID == resourceID {
			return resource, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateResource replaces a Vampire's existing resource with the new one
// based on the new resource's ID.
func (v *OldVampire) UpdateResource(newResource *Resource) error {
	oldResource, err := v.FindResource(newResource.ID)
	if err != nil {
		return err
	}

	*oldResource = *newResource

	return nil
}

// AddCharacter adds a character to the Vampire.
func (v *OldVampire) AddCharacter(character *Character) {
	character.ID = len(v.Characters) + 1
	v.Characters = append(v.Characters, character)
}

// FindCharacter retrieves a character based on an ID from the Vampire's list of
// characters.
func (v *OldVampire) FindCharacter(characterID int) (*Character, error) {
	for _, character := range v.Characters {
		if character.ID == characterID {
			return character, nil
		}
	}

	return nil, ErrNotFound
}

// UpdateCharacter replaces a Vampire's existing character with the new one
// based on the new character's ID.
func (v *OldVampire) UpdateCharacter(newCharacter *Character) error {
	oldCharacter, err := v.FindCharacter(newCharacter.ID)
	if err != nil {
		return err
	}

	*oldCharacter = *newCharacter

	return nil
}

// AddDescriptor adds a descriptor to the indicated character.
func (v *OldVampire) AddDescriptor(characterID int, descriptor string) error {
	character, err := v.FindCharacter(characterID)
	if err != nil {
		return err
	}

	character.AddDescriptor(descriptor)

	return nil
}

// AddMark adds a mark to the Vampire.
func (v *OldVampire) AddMark(mark *Mark) {
	mark.ID = len(v.Marks) + 1
	v.Marks = append(v.Marks, mark)
}
