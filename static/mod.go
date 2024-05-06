package static

import "embed"

//go:embed *.css *.js *.ico
var FS embed.FS
