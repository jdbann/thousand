package app

import (
	"io"
	"net/http"
	"strings"
)

var content = `
<html>
<head>
<title>Thousand</title>
</head>
<body>
<h1>Thousand</h1>
</body>
</html>
`

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
}

// Routes returns the configured set of routes for the app to be used by an HTTP
// server.
func (app *App) Routes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, strings.TrimSpace(content))
	})
}
