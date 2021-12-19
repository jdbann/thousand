package templates

import (
	"net/http"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/models"
)

func (r *Renderer) NewCharacter(w http.ResponseWriter, req *http.Request, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, req, "characters/new", data)
}

func (r *Renderer) NewExperience(w http.ResponseWriter, req *http.Request, m models.Memory) error {
	data := map[string]interface{}{
		"memory": m,
	}

	return r.render(w, req, "experiences/new", data)
}

func (r *Renderer) NewMark(w http.ResponseWriter, req *http.Request, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, req, "marks/new", data)
}

func (r *Renderer) NewSkill(w http.ResponseWriter, req *http.Request, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, req, "skills/new", data)
}

func (r *Renderer) NewResource(w http.ResponseWriter, req *http.Request, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, req, "resources/new", data)
}

func (r *Renderer) NewUser(w http.ResponseWriter, req *http.Request, f *form.NewUserForm) error {
	data := map[string]interface{}{
		"form": f,
	}

	return r.render(w, req, "users/new", data)
}

func (r *Renderer) ShowVampires(w http.ResponseWriter, req *http.Request, v []models.Vampire) error {
	data := map[string]interface{}{
		"vampires": v,
	}

	return r.render(w, req, "vampires/index", data)
}

func (r *Renderer) NewVampire(w http.ResponseWriter, req *http.Request) error {
	return r.render(w, req, "vampires/new", map[string]interface{}{})
}

func (r *Renderer) ShowVampire(w http.ResponseWriter, req *http.Request, v models.Vampire) error {
	data := map[string]interface{}{
		"vampire": v,
	}

	return r.render(w, req, "vampires/show", data)
}
