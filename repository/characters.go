package repository

import (
	"context"
	"errors"
	"fmt"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func (m *Repository) CreateCharacter(ctx context.Context, vampireID uuid.UUID, params models.CreateCharacterParams) (models.Character, error) {
	var characterType queries.CharacterType
	switch params.Type {
	case "mortal":
		characterType = queries.CharacterTypeMortal
	case "immortal":
		characterType = queries.CharacterTypeImmortal
	default:
		return models.Character{}, fmt.Errorf("unrecognised character type: %q", params.Type)
	}

	dbParams := queries.CreateCharacterParams{
		VampireID: vampireID,
		Name:      params.Name,
		Type:      characterType,
	}

	dbCharacter, err := m.queries.CreateCharacter(ctx, dbParams)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.ForeignKeyViolation && pgErr.ConstraintName == "characters_vampire_id_fkey" {
			return models.Character{}, models.ErrNotFound.Cause(err)
		}

		return models.Character{}, err
	} else if err != nil {
		return models.Character{}, err
	}

	return newCharacter(dbCharacter), nil
}
