package models

import (
	"context"
	"errors"
	"fmt"

	"emailaddress.horse/thousand/db"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Models holds a configured *db.Queries which can be used to load values from
// the DB and return them as model types from this package.
type Models struct {
	*db.Queries
}

// DBTX is an alias of db.DBTX to ensure that the models package completely
// encapsulates any DB interactions.
type DBTX = db.DBTX

// NewModels returns a newly configured models.Models struct.
func NewModels(dbtx DBTX) *Models {
	return &Models{db.New(dbtx)}
}

// CreateVampire attempts to create a new vampire in the DB with the provided
// name.
func (m *Models) CreateVampire(ctx context.Context, name string) (Vampire, error) {
	v, err := m.Queries.CreateVampire(ctx, name)
	if err != nil {
		return Vampire{}, err
	}

	var createMemoriesParams = make([]uuid.UUID, vampireMemorySize)
	for i := range createMemoriesParams {
		createMemoriesParams[i] = v.ID
	}

	dbMemories, err := m.Queries.CreateMemories(ctx, createMemoriesParams)
	if err != nil {
		return Vampire{}, err
	}

	memories := make([]Memory, len(dbMemories))
	for i, dbMemory := range dbMemories {
		memories[i] = newMemory(dbMemory, make([]db.Experience, 0, 3))
	}

	return newVampire(v, memories, []Skill{}, []Resource{}, []Character{}, []Mark{}), nil
}

// GetVampire attempts to retrieve a vampire from the DB with the provided ID.
func (m *Models) GetVampire(ctx context.Context, id uuid.UUID) (Vampire, error) {
	v, err := m.Queries.GetVampire(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Vampire{}, ErrNotFound.Cause(err)
	} else if err != nil {
		return Vampire{}, err
	}

	dbMemories, err := m.Queries.GetMemoriesForVampire(ctx, id)
	if err != nil {
		return Vampire{}, err
	}

	dbExperiences, err := m.Queries.GetExperiencesForVampire(ctx, id)
	if err != nil {
		return Vampire{}, err
	}

	memories := make([]Memory, len(dbMemories))
	for i, dbMemory := range dbMemories {
		experiences := make([]db.Experience, 0, 3)

		for _, experience := range dbExperiences {
			if experience.MemoryID == dbMemory.ID {
				experiences = append(experiences, experience)
			}
		}

		memories[i] = newMemory(dbMemory, experiences)
	}

	dbSkills, err := m.Queries.GetSkillsForVampire(ctx, id)
	if err != nil {
		return Vampire{}, err
	}

	skills := make([]Skill, len(dbSkills))
	for i, dbSkill := range dbSkills {
		skills[i] = newSkill(dbSkill)
	}

	dbResources, err := m.Queries.GetResourcesForVampire(ctx, id)
	if err != nil {
		return Vampire{}, err
	}

	resources := make([]Resource, len(dbResources))
	for i, dbResource := range dbResources {
		resources[i] = newResource(dbResource)
	}

	dbCharacters, err := m.Queries.GetCharactersForVampire(ctx, id)
	if err != nil {
		return Vampire{}, err
	}

	characters := make([]Character, len(dbCharacters))
	for i, dbCharacter := range dbCharacters {
		characters[i] = newCharacter(dbCharacter)
	}

	dbMarks, err := m.Queries.GetMarksForVampire(ctx, id)
	if err != nil {
		return Vampire{}, err
	}

	marks := make([]Mark, len(dbMarks))
	for i, dbMark := range dbMarks {
		marks[i] = newMark(dbMark)
	}

	return newVampire(v, memories, skills, resources, characters, marks), nil
}

// GetVampires attempts to retrieve all the vampires from the DB.
func (m *Models) GetVampires(ctx context.Context) ([]Vampire, error) {
	vs, err := m.Queries.GetVampires(ctx)
	if err != nil {
		return []Vampire{}, err
	}

	nvs := make([]Vampire, len(vs))
	for i, v := range vs {
		nvs[i] = newVampire(v, []Memory{}, []Skill{}, []Resource{}, []Character{}, []Mark{})
	}

	return nvs, nil
}

func (m *Models) GetMemory(ctx context.Context, vampireID, id uuid.UUID) (Memory, error) {
	params := db.GetMemoryParams{
		VampireID: vampireID,
		MemoryID:  id,
	}

	dbMemory, err := m.Queries.GetMemory(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return Memory{}, ErrNotFound.Cause(err)
	} else if err != nil {
		return Memory{}, err
	}

	// TODO: Also return experiences?
	return newMemory(dbMemory, []db.Experience{}), nil
}

// AddExperience attempts to add a new experience to the DB for the provided
// memory.
func (m *Models) AddExperience(ctx context.Context, vampireID, memoryID uuid.UUID, description string) (Experience, error) {
	params := db.CreateExperienceParams{
		VampireID:   vampireID,
		MemoryID:    memoryID,
		Description: description,
	}

	dbExperience, err := m.Queries.CreateExperience(ctx, params)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == PgErrCodeMemoryFull {
		return Experience{}, ErrMemoryFull.Cause(err)
	} else if err != nil {
		return Experience{}, err
	}

	return newExperience(dbExperience), nil
}

// GetExperiences attempts to retrieve all the experiences from the DB for the
// provided vampire.
func (m *Models) GetExperiences(ctx context.Context, vampireID uuid.UUID) ([]Experience, error) {
	dbExperiences, err := m.Queries.GetExperiencesForVampire(ctx, vampireID)
	if err != nil {
		return nil, err
	}

	var experiences []Experience
	for _, dbExperience := range dbExperiences {
		experiences = append(experiences, newExperience(dbExperience))
	}

	return experiences, nil
}

// AddSkill attempts to add a new skill to the DB for the provided vampire.
func (m *Models) AddSkill(ctx context.Context, vampireID uuid.UUID, description string) (Skill, error) {
	params := db.CreateSkillParams{
		VampireID:   vampireID,
		Description: description,
	}

	dbSkill, err := m.Queries.CreateSkill(ctx, params)
	if err != nil {
		return Skill{}, err
	}

	return newSkill(dbSkill), nil
}

type AddResourceParams struct {
	Description string `form:"description"`
	Stationary  bool   `form:"stationary"`
}

// AddResource attempts to add a new resource to the DB for the provided vampire.
func (m *Models) AddResource(ctx context.Context, vampireID uuid.UUID, params AddResourceParams) (Resource, error) {
	dbParams := db.CreateResourceParams{
		VampireID:   vampireID,
		Description: params.Description,
		Stationary:  params.Stationary,
	}

	dbResource, err := m.Queries.CreateResource(ctx, dbParams)
	if err != nil {
		return Resource{}, err
	}

	return newResource(dbResource), nil
}

type AddCharacterParams struct {
	Name string `form:"name"`
	Type string `form:"type"`
}

func (m *Models) AddCharacter(ctx context.Context, vampireID uuid.UUID, params AddCharacterParams) (Character, error) {
	var characterType db.CharacterType
	switch params.Type {
	case "mortal":
		characterType = db.CharacterTypeMortal
	case "immortal":
		characterType = db.CharacterTypeImmortal
	default:
		return Character{}, fmt.Errorf("unrecognised character type: %q", params.Type)
	}

	dbParams := db.CreateCharacterParams{
		VampireID: vampireID,
		Name:      params.Name,
		Type:      characterType,
	}

	dbCharacter, err := m.Queries.CreateCharacter(ctx, dbParams)
	if err != nil {
		return Character{}, err
	}

	return newCharacter(dbCharacter), nil
}

func (m *Models) AddMark(ctx context.Context, vampireID uuid.UUID, description string) (Mark, error) {
	params := db.CreateMarkParams{
		VampireID:   vampireID,
		Description: description,
	}

	dbMark, err := m.Queries.CreateMark(ctx, params)
	if err != nil {
		return Mark{}, err
	}

	return newMark(dbMark), nil
}
