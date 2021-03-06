package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"emailaddress.horse/thousand/middleware"
	"emailaddress.horse/thousand/session"
)

//go:embed views
var viewTemplates embed.FS

//go:embed layouts/*.tmpl
var layoutTemplates embed.FS

type Renderer struct {
	store       *session.Store
	templateMap map[string]*template.Template
}

type RendererOptions struct {
	Store *session.Store
}

func NewRenderer(opts RendererOptions) *Renderer {
	templateMap := map[string]*template.Template{}

	// Traverse the views directory, descending into each folder
	err := fs.WalkDir(viewTemplates, "views", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		// Don't take further action if the entry is a folder or is not a template
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Get the contents of the file
		viewBytes, err := viewTemplates.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		// Build the name we will use for the template by removing the views/ prefix
		// and .tmpl suffix
		name := strings.TrimPrefix(path, "views/")
		name = strings.TrimSuffix(name, ".tmpl")

		// Setup the new template for the view
		viewTemplate := template.New(name)

		// Add the helpers
		viewTemplate.Funcs(helpers)

		// Parse the contents of the file
		viewTemplate, err = viewTemplate.Parse(string(viewBytes))
		if err != nil {
			log.Fatal(err)
		}

		// Also parse the contents of the layouts directory to make sure the view
		// has access to all layouts
		viewTemplate, err = viewTemplate.ParseFS(layoutTemplates, "layouts/*.tmpl")
		if err != nil {
			log.Fatal(err)
		}

		// Add the view to the templates map
		templateMap[name] = viewTemplate

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Renderer{
		store:       opts.Store,
		templateMap: templateMap,
	}
}

func (r *Renderer) render(w http.ResponseWriter, req *http.Request, name string, data map[string]interface{}) error {
	flashes, err := r.store.GetFlashes(req, w)
	if err != nil {
		return err
	}
	data["flashes"] = flashes

	currentUser, ok := middleware.MaybeCurrentUser(req.Context())
	if ok {
		data["currentUser"] = currentUser
	}

	view, ok := r.templateMap[name]
	if !ok {
		return fmt.Errorf("No template found with name: %q", name)
	}

	return view.ExecuteTemplate(w, name, data)
}
