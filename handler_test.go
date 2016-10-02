package main

import (
	"net/http"
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
	rr := CallAndRecord(t, config, "/")
	assert.Equal(t, rr.Code, http.StatusOK)
	AssertHTML(t, rr.Body.String(), `
<!DOCTYPE html>
<html>
  <body>
    <ul>
      <li>
        thriftrw - github.com/thriftrw/thriftrw-go
      </li>
      <li>
        yarpc - github.com/yarpc/yarpc-go
      </li>
    </ul>
  </body>
</html>
`)
}

func TestPackageShouldExist(t *testing.T) {
	rr := CallAndRecord(t, config, "/yarpc")
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), `
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
	rr := CallAndRecord(t, config, "/nonexistent")
	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestTrailingSlash(t *testing.T) {
	rr := CallAndRecord(t, config, "/yarpc/")
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), `
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
	rr := CallAndRecord(t, config, "/yarpc/heeheehee")
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), `
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

	rr = CallAndRecord(t, config, "/yarpc/heehee/hawhaw")
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), `
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
