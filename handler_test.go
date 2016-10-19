package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config = `

url: go.uber.org
packages:
  yarpc:
    repo: github.com/yarpc/yarpc-go
  thriftrw:
    repo: github.com/thriftrw/thriftrw-go

`

func TestIndex(t *testing.T) {
	rr := CallAndRecord(t, config, "GET", "/")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <body>
        <ul>
            <li>thriftrw - github.com/thriftrw/thriftrw-go</li>
            <li>yarpc - github.com/yarpc/yarpc-go</li>
        </ul>
    </body>
</html>
`)
}

func TestPackageShouldExist(t *testing.T) {
	rr := CallAndRecord(t, config, "GET", "/yarpc")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://godoc.org/go.uber.org/yarpc">
    </head>
    <body>
        Nothing to see here. Please <a href="https://godoc.org/go.uber.org/yarpc">move along</a>.
    </body>
</html>
`)
}

func TestNonExistentPackageShould404(t *testing.T) {
	rr := CallAndRecord(t, config, "GET", "/nonexistent")
	AssertResponse(t, rr, 404, `
404 page not found
`)
	assert.Equal(t, "no-cache", rr.Header().Get("Cache-Control"))
}

func TestTrailingSlash(t *testing.T) {
	rr := CallAndRecord(t, config, "GET", "/yarpc/")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://godoc.org/go.uber.org/yarpc/">
    </head>
    <body>
        Nothing to see here. Please <a href="https://godoc.org/go.uber.org/yarpc/">move along</a>.
    </body>
</html>
`)
}

func TestDeepImports(t *testing.T) {
	rr := CallAndRecord(t, config, "GET", "/yarpc/heeheehee")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://godoc.org/go.uber.org/yarpc/heeheehee">
    </head>
    <body>
        Nothing to see here. Please <a href="https://godoc.org/go.uber.org/yarpc/heeheehee">move along</a>.
    </body>
</html>
`)

	rr = CallAndRecord(t, config, "GET", "/yarpc/heehee/hawhaw")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta name="go-source" content="go.uber.org/yarpc https://github.com/yarpc/yarpc-go https://github.com/yarpc/yarpc-go/tree/master{/dir} https://github.com/yarpc/yarpc-go/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://godoc.org/go.uber.org/yarpc/heehee/hawhaw">
    </head>
    <body>
        Nothing to see here. Please <a href="https://godoc.org/go.uber.org/yarpc/heehee/hawhaw">move along</a>.
    </body>
</html>
`)
}

func TestMethodNotAllowed(t *testing.T) {
	methods := []string{"POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	uris := []string{"/", "/yarpc"}
	for _, method := range methods {
		for _, uri := range uris {
			t.Run(fmt.Sprintf("%s => %s", method, uri), func(t *testing.T) {
				rr := CallAndRecord(t, config, method, uri)
				AssertResponse(t, rr, 405, "\n405 method not allowed\n")
				assert.Equal(t, "no-cache", rr.Header().Get("Cache-Control"))
			})
		}
	}
}

func TestInternalServerError(t *testing.T) {
	rr := CallAndRecord(t, config, "GET", "/panic")
	AssertResponse(t, rr, 500, "\n500 internal server error\n")
	assert.Equal(t, "no-cache", rr.Header().Get("Cache-Control"))
}
