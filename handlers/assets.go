package handlers

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func Assets(r chi.Router, assetsFS fs.FS) {
	assets := indexHidingFS{
		fs: http.FS(assetsFS),
	}

	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(assets))

		rw := newNotFoundRescuer(w, http.StatusNotFound)
		fs.ServeHTTP(rw, r)
		if rw.rescued {
			handleError(w, NotFoundError)
		}
	})
}

type indexHidingFS struct {
	fs http.FileSystem
}

func (fsys indexHidingFS) Open(name string) (http.File, error) {
	file, err := fsys.fs.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fs.ErrNotExist
	}

	return file, nil
}
