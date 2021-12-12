package server

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/go-chi/chi/v5"
)

type route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func PrintRoutes(r chi.Router) error {
	var methodOrder = map[string]int{
		"GET":     0,
		"POST":    1,
		"PUT":     2,
		"PATCH":   3,
		"DELETE":  4,
		"HEAD":    5,
		"CONNECT": 6,
		"OPTIONS": 7,
		"TRACE":   8,
	}

	routes := []*route{}

	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, &route{
			Method: method,
			Path:   path,
		})
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		return err
	}

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return methodOrder[routes[i].Method] < methodOrder[routes[j].Method]
		}

		return routes[i].Path < routes[j].Path
	})

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(writer, "Method\tPath\n")
	for _, route := range routes {
		fmt.Fprintf(writer, "%s\t%s\n", route.Method, route.Path)
	}
	return writer.Flush()
}
