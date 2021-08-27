package static

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed  css
var _assets embed.FS

// Middleware returns an echo.MiddlewareFunc which will handle serving static
// assets.
func Middleware() echo.MiddlewareFunc {
	config := middleware.DefaultStaticConfig

	config.Filesystem = http.FS(_assets)

	return middleware.StaticWithConfig(config)
}
