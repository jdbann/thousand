package models

import (
	"strings"

	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
)

type Character struct {
	ID   uuid.UUID
	Name string
	Type string
}

func newCharacter(dbCharacter queries.Character) Character {
	return Character{
		ID:   dbCharacter.ID,
		Name: dbCharacter.Name,
		Type: strings.Title(string(dbCharacter.Type)),
	}
}
