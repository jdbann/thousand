package handlers

import (
	"context"

	"emailaddress.horse/thousand/models"
	"github.com/google/uuid"
)

type characterCreator interface {
	CreateCharacter(context.Context, uuid.UUID, models.CreateCharacterParams) (models.Character, error)
}

type experienceCreator interface {
	CreateExperience(context.Context, uuid.UUID, uuid.UUID, string) (models.Experience, error)
}

type markCreator interface {
	CreateMark(context.Context, uuid.UUID, string) (models.Mark, error)
}

type memoryGetter interface {
	GetMemory(context.Context, uuid.UUID, uuid.UUID) (models.Memory, error)
}

type resourceCreator interface {
	CreateResource(context.Context, uuid.UUID, models.CreateResourceParams) (models.Resource, error)
}

type skillCreator interface {
	CreateSkill(context.Context, uuid.UUID, string) (models.Skill, error)
}

type vampireCreator interface {
	CreateVampire(context.Context, uuid.UUID, string) (models.Vampire, error)
}

type vampireGetter interface {
	GetVampire(context.Context, uuid.UUID) (models.Vampire, error)
}

type vampiresGetter interface {
	GetVampires(context.Context) ([]models.Vampire, error)
}
