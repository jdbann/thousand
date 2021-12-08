package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"strings"
)

type Renderer struct {
	templateMap map[string]*template.Template
}

func NewRenderer() *Renderer {
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
		templateMap: templateMap,
	}
}

func (r *Renderer) render(w io.Writer, name string, data interface{}) error {
	view, ok := r.templateMap[name]
	if !ok {
		return fmt.Errorf("No template found with name: %q", name)
	}

	return view.ExecuteTemplate(w, name, data)
}
