// Package templates exposes the template used by Sally
// to render the HTML pages.
package templates

import _ "embed" // needed for go:embed

// Index holds the contents of the index.html template.
//
//go:embed index.html
var Index string

// Package holds the contents of the package.html template.
//
//go:embed package.html
var Package string
