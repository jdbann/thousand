package templates

import (
	"fmt"
	"html/template"

	"github.com/google/uuid"
)

var helpers = template.FuncMap{
	"newCharacterPath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/characters/new", vampireID)
	},
	"createCharacterPath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/characters", vampireID)
	},

	"newExperiencePath": func(vampireID, memoryID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/memories/%s/experiences/new", vampireID, memoryID)
	},
	"createExperiencePath": func(vampireID, memoryID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/memories/%s/experiences", vampireID, memoryID)
	},

	"newMarkPath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/marks/new", vampireID)
	},
	"createMarkPath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/marks", vampireID)
	},

	"newResourcePath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/resources/new", vampireID)
	},
	"createResourcePath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/resources", vampireID)
	},

	"sessionPath": func() string {
		return "/session"
	},
	"newSessionPath": func() string {
		return "/session/new"
	},

	"newSkillPath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/skills/new", vampireID)
	},
	"createSkillPath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s/skills", vampireID)
	},

	"userPath": func() string {
		return "/user"
	},
	"newUserPath": func() string {
		return "/user/new"
	},

	"vampiresPath": func() string {
		return "/vampires"
	},
	"newVampirePath": func() string {
		return "/vampires/new"
	},
	"vampirePath": func(vampireID uuid.UUID) string {
		return fmt.Sprintf("/vampires/%s", vampireID)
	},
}
