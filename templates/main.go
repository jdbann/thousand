package templates

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed *.tmpl
var _templates embed.FS

// NewRenderer returns a value that implements the echo.Renderer interface for
// rendering the HTML templates required by this application.
func NewRenderer() echo.Renderer {
	templates, err := template.ParseFS(_templates, "*.tmpl")
	if err != nil {
		panic(err)
	}

	return &renderer{templates}
}

type renderer struct {
	templates *template.Template
}

// Render allows the renderer type to adhere to the echo.Renderer interface for
// rendering HTML templates.
func (t *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name+".tmpl", data)
}
