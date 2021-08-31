package templates

import _ "embed"

//go:embed index.html
var Index string

//go:embed package.html
var Package string
