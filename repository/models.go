package repository

import (
	"strings"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/repository/queries"
)

func newCharacter(dbCharacter queries.Character) models.Character {
	return models.Character{
		ID:   dbCharacter.ID,
		Name: dbCharacter.Name,
		Type: strings.Title(string(dbCharacter.Type)),
	}
}

func newExperience(experience queries.Experience) models.Experience {
	return models.Experience{
		ID:          experience.ID,
		MemoryID:    experience.MemoryID,
		Description: experience.Description,
	}
}

func newMark(dbMark queries.Mark) models.Mark {
	return models.Mark{
		ID:          dbMark.ID,
		Description: dbMark.Description,
	}
}

func newMemory(dbMemory queries.Memory, dbExperiences []queries.Experience) models.Memory {
	var experiences = make([]models.Experience, len(dbExperiences))
	for i, dbExperience := range dbExperiences {
		experiences[i] = newExperience(dbExperience)
	}

	return models.Memory{
		ID:          dbMemory.ID,
		VampireID:   dbMemory.VampireID,
		Experiences: experiences,
	}
}

func newResource(dbResource queries.Resource) models.Resource {
	return models.Resource{
		ID:          dbResource.ID,
		VampireID:   dbResource.VampireID,
		Description: dbResource.Description,
		Stationary:  dbResource.Stationary,
	}
}

func newSkill(dbSkill queries.Skill) models.Skill {
	return models.Skill{
		ID:          dbSkill.ID,
		VampireID:   dbSkill.VampireID,
		Description: dbSkill.Description,
	}
}

func newVampire(dbVampire queries.Vampire, memories []models.Memory, skills []models.Skill, resources []models.Resource, characters []models.Character, marks []models.Mark) models.Vampire {
	return models.Vampire{
		ID:         dbVampire.ID,
		Name:       dbVampire.Name,
		Memories:   memories,
		Skills:     skills,
		Resources:  resources,
		Characters: characters,
		Marks:      marks,
	}
}
