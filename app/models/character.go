package models

import (
	"fmt"
	"strings"

	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
)

type Character struct {
	ID   uuid.UUID
	Name string
	Type string
}

func newCharacter(dbCharacter db.Character) Character {
	return Character{
		ID:   dbCharacter.ID,
		Name: dbCharacter.Name,
		Type: strings.Title(string(dbCharacter.Type)),
	}
}

// OldCharacter holds the details of a character known by a Vampire.
type OldCharacter struct {
	ID          int
	Name        string   `form:"name"`
	Descriptors []string `form:"descriptor[]"`
	Type        string   `form:"type"`
	Deceased    bool     `form:"deceased"`
}

// Description returns a string joining the name, descriptors and type of a
// Character.
func (c *OldCharacter) Description() string {
	components := append([]string{c.Name}, c.Descriptors...)

	description := fmt.Sprintf("%s.", strings.Join(components, ", "))

	if c.Type != "" {
		description = fmt.Sprintf("%s (%s)", description, strings.Title(c.Type))
	}

	return description
}

// AddDescriptor appends a new descriptor to the Character's current list of
// descriptors.
func (c *OldCharacter) AddDescriptor(descriptor string) {
	c.Descriptors = append(c.Descriptors, descriptor)
}
