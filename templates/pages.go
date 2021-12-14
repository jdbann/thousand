package templates

import (
	"net/http"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/models"
)

func (r *Renderer) NewCharacter(w http.ResponseWriter, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, "characters/new", data)
}

func (r *Renderer) NewExperience(w http.ResponseWriter, m models.Memory) error {
	data := map[string]interface{}{
		"memory": m,
	}

	return r.render(w, "experiences/new", data)
}

func (r *Renderer) NewMark(w http.ResponseWriter, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, "marks/new", data)
}

func (r *Renderer) NewSkill(w http.ResponseWriter, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, "skills/new", data)
}

func (r *Renderer) NewResource(w http.ResponseWriter, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, "resources/new", data)
}

func (r *Renderer) NewUser(w http.ResponseWriter, f *form.NewUserForm) error {
	data := map[string]interface{}{
		"form": f,
	}

	return r.render(w, "users/new", data)
}

func (r *Renderer) ShowVampires(w http.ResponseWriter, v []models.Vampire) error {
	data := map[string]interface{}{
		"vampires": v,
	}

	return r.render(w, "vampires/index", data)
}

func (r *Renderer) NewVampire(w http.ResponseWriter) error {
	return r.render(w, "vampires/new", nil)
}

func (r *Renderer) ShowVampire(w http.ResponseWriter, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, "vampires/show", data)
}
