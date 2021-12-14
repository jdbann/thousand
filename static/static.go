package static

import (
	"embed"
)

//nolint:typecheck
//go:embed  css js
var Assets embed.FS
