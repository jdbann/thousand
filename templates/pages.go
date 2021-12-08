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
