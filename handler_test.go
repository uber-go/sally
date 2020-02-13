package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var config = `

url: go.uber.org
packages:
  thriftrw:
    repo: github.com/thriftrw/thriftrw-go
  yarpc:
    repo: github.com/yarpc/yarpc-go

`

func TestIndex(t *testing.T) {
	rr := CallAndRecord(t, config, "/")
	assert.Equal(t, 200, rr.Code)

	body := rr.Body.String()
	assert.Contains(t, body, "github.com/thriftrw/thriftrw-go")
	assert.Contains(t, body, "github.com/yarpc/yarpc-go")
}

func TestPackageShouldExist(t *testing.T) {
	rr := CallAndRecord(t, config, "/yarpc")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc">
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc">move along</a>.
    </body>
</html>
`)
}

func TestNonExistentPackageShould404(t *testing.T) {
	rr := CallAndRecord(t, config, "/nonexistent")
	AssertResponse(t, rr, 404, `
404 page not found
`)
}

func TestTrailingSlash(t *testing.T) {
	rr := CallAndRecord(t, config, "/yarpc/")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc/">
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc/">move along</a>.
    </body>
</html>
`)
}

func TestDeepImports(t *testing.T) {
	rr := CallAndRecord(t, config, "/yarpc/heeheehee")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc/heeheehee">
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc/heeheehee">move along</a>.
    </body>
</html>
`)

	rr = CallAndRecord(t, config, "/yarpc/heehee/hawhaw")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc/heehee/hawhaw">
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc/heehee/hawhaw">move along</a>.
    </body>
</html>
`)
}
