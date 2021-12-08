package templates

import (
	"net/http"

	"emailaddress.horse/thousand/models"
)

func (r *Renderer) ShowVampires(w http.ResponseWriter, v []models.Vampire) error {
	data := map[string]interface{}{
		"vampires": v,
	}

	return r.render(w, "vampires/index", data)
}

func (r *Renderer) NewVampire(w http.ResponseWriter) error {
	return r.render(w, "vampires/new", nil)
}
