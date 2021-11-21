package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed views
var viewTemplates embed.FS

//go:embed layouts/*.tmpl
var layoutTemplates embed.FS

type routeReverser interface {
	Reverse(name string, params ...interface{}) string
}

// NewRenderer returns a value that implements the echo.Renderer interface for
// rendering the HTML templates required by this application.
//
// It creates templates for each file in templates/views using the path and
// filename without the file extension. So "templates/views/vampires/show.tmpl"
// has the name "vampires/show".
func NewRenderer(rr routeReverser) echo.Renderer {
	views := map[string]*template.Template{}

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

		// Add the Reverse function to the template
		viewTemplate.Funcs(template.FuncMap{
			"reverse": rr.Reverse,
		})

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
		views[name] = viewTemplate

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return &renderer{views}
}

type renderer struct {
	views map[string]*template.Template
}

// Render allows the renderer type to adhere to the echo.Renderer interface for
// rendering HTML templates.
func (t *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	view, ok := t.views[name]
	if !ok {
		return fmt.Errorf("No template found with name: %q", name)
	}

	return view.ExecuteTemplate(w, name, data)
}
