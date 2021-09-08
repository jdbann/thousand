package static

import (
	"crypto/sha1"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	//go:embed  css
	_assets embed.FS

	etagMap map[string]string = make(map[string]string)
)

func init() {
	// Traverse the _assets file system to precompile strong ETag values for each
	// asset.
	err := fs.WalkDir(_assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		// Don't take further action if the entry is a folder
		if d.IsDir() {
			return nil
		}

		// Get the contents of the file
		assetBytes, err := _assets.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		// Build an ETag for the file and add it to the map
		hash := sha1.Sum(assetBytes)
		etag := fmt.Sprintf("\"%x\"", hash)
		etagMap[path] = etag

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

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
// requested file from the underlying http.FileSystem. It also adds an ETag
// header based on the precompiled map of strong ETag values built in init().
func (efs etagFS) Open(name string) (http.File, error) {
	etag, ok := etagMap[name]
	if ok {
		efs.context.Response().Header().Set("ETag", etag)
	}

	return efs.realFS.Open(name)
}
