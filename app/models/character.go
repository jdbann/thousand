package models

import (
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
