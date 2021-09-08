package static

import (
	"bytes"
	"crypto/sha1"
	"embed"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed  css
var _assets embed.FS

// Middleware returns an echo.MiddlewareFunc which will handle serving static
// assets.
func Middleware() echo.MiddlewareFunc {
	// Echo's middleware.Static doesn't offer a mechanism for setting an ETag
	// prior to using http.ServeContent based on the requested file. To cheat
	// this, I am wrapping the embed.FS I use for the assets as an etagFS type.
	// This type also requires the echo.Context in order to write the appropriate
	// headers - so it must be built inside an echo.HandlerFunc.
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			config := middleware.DefaultStaticConfig
			config.Filesystem = etagFS{c, http.FS(_assets)}

			// Building a new middleware.StaticWithConfig each time so this could
			// definitely be improved - but it saves duplicating the whole package for
			// the sake of adding ETag headers.
			return middleware.StaticWithConfig(config)(next)(c)
		}
	}
}

type etagFS struct {
	context echo.Context
	realFS  http.FileSystem
}

// Open satisfies the http.FileSystem interface for etagFS, opening the
// requested file from the underlying http.FileSystem. It also calculates a hash
// of the opened file and adds it as a weak ETag header.
func (efs etagFS) Open(name string) (http.File, error) {
	file, err := efs.realFS.Open(name)
	if err != nil {
		return nil, err
	}

	// Hashes are calculated for each request. This could be improved by building
	// a map of file names to hashes on init() of the package and looking up the
	// hash instead.
	var contents bytes.Buffer
	if _, err := contents.ReadFrom(file); err != nil {
		return nil, err
	}
	hash := sha1.Sum(contents.Bytes())

	etag := fmt.Sprintf("W/\"%x\"", hash)
	efs.context.Response().Header().Set("ETag", etag)

	return file, nil
}
