package repository

import (
	"context"
	"errors"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// CreateVampire attempts to create a new vampire in the DB with the provided
// name.
func (m *Repository) CreateVampire(ctx context.Context, userID uuid.UUID, name string) (models.Vampire, error) {
	params := queries.CreateVampireParams{
		Name:   name,
		UserID: userID,
	}

	v, err := m.queries.CreateVampire(ctx, params)
	if err != nil {
		return models.Vampire{}, err
	}

	var createMemoriesParams = make([]uuid.UUID, models.VampireMemorySize)
	for i := range createMemoriesParams {
		createMemoriesParams[i] = v.ID
	}

	dbMemories, err := m.queries.CreateMemories(ctx, createMemoriesParams)
	if err != nil {
		return models.Vampire{}, err
	}

	memories := make([]models.Memory, len(dbMemories))
	for i, dbMemory := range dbMemories {
		memories[i] = newMemory(dbMemory, make([]queries.Experience, 0, 3))
	}

	return newVampire(v, memories, []models.Skill{}, []models.Resource{}, []models.Character{}, []models.Mark{}), nil
}

// GetVampire attempts to retrieve a vampire from the DB with the provided ID.
func (m *Repository) GetVampire(ctx context.Context, id uuid.UUID) (models.Vampire, error) {
	v, err := m.queries.GetVampire(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Vampire{}, models.ErrNotFound.Cause(err)
	} else if err != nil {
		return models.Vampire{}, err
	}

	dbMemories, err := m.queries.GetMemoriesForVampire(ctx, id)
	if err != nil {
		return models.Vampire{}, err
	}

	dbExperiences, err := m.queries.GetExperiencesForVampire(ctx, id)
	if err != nil {
		return models.Vampire{}, err
	}

	memories := make([]models.Memory, len(dbMemories))
	for i, dbMemory := range dbMemories {
		experiences := make([]queries.Experience, 0, 3)

		for _, experience := range dbExperiences {
			if experience.MemoryID == dbMemory.ID {
				experiences = append(experiences, experience)
			}
		}

		memories[i] = newMemory(dbMemory, experiences)
	}

	dbSkills, err := m.queries.GetSkillsForVampire(ctx, id)
	if err != nil {
		return models.Vampire{}, err
	}

	skills := make([]models.Skill, len(dbSkills))
	for i, dbSkill := range dbSkills {
		skills[i] = newSkill(dbSkill)
	}

	dbResources, err := m.queries.GetResourcesForVampire(ctx, id)
	if err != nil {
		return models.Vampire{}, err
	}

	resources := make([]models.Resource, len(dbResources))
	for i, dbResource := range dbResources {
		resources[i] = newResource(dbResource)
	}

	dbCharacters, err := m.queries.GetCharactersForVampire(ctx, id)
	if err != nil {
		return models.Vampire{}, err
	}

	characters := make([]models.Character, len(dbCharacters))
	for i, dbCharacter := range dbCharacters {
		characters[i] = newCharacter(dbCharacter)
	}

	dbMarks, err := m.queries.GetMarksForVampire(ctx, id)
	if err != nil {
		return models.Vampire{}, err
	}

	marks := make([]models.Mark, len(dbMarks))
	for i, dbMark := range dbMarks {
		marks[i] = newMark(dbMark)
	}

	return newVampire(v, memories, skills, resources, characters, marks), nil
}

// GetVampires attempts to retrieve all the vampires from the DB.
func (m *Repository) GetVampires(ctx context.Context) ([]models.Vampire, error) {
	vs, err := m.queries.GetVampires(ctx)
	if err != nil {
		return []models.Vampire{}, err
	}

	nvs := make([]models.Vampire, len(vs))
	for i, v := range vs {
		nvs[i] = newVampire(v, []models.Memory{}, []models.Skill{}, []models.Resource{}, []models.Character{}, []models.Mark{})
	}

	return nvs, nil
}
