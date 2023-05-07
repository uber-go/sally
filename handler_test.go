package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var config = `

url: go.uber.org
packages:
  thriftrw:
    repo: github.com/thriftrw/thriftrw-go
  yarpc:
    repo: github.com/yarpc/yarpc-go
  zap:
    url: go.uberalt.org
    repo: github.com/uber-go/zap
    description: A fast, structured logging library.

`

func TestIndex(t *testing.T) {
	rr := CallAndRecord(t, config, "/")
	assert.Equal(t, 200, rr.Code)

	body := rr.Body.String()
	assert.Contains(t, body, "github.com/thriftrw/thriftrw-go")
	assert.Contains(t, body, "github.com/yarpc/yarpc-go")
	assert.Contains(t, body, "A fast, structured logging library.")
}

func TestPackageShouldExist(t *testing.T) {
	rr := CallAndRecord(t, config, "/yarpc")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
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
	<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc/heehee/hawhaw">
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc/heehee/hawhaw">move along</a>.
    </body>
</html>
`)
}

func TestPackageLevelURL(t *testing.T) {
	rr := CallAndRecord(t, config, "/zap")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uberalt.org/zap git https://github.com/uber-go/zap">
	<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uberalt.org/zap">
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uberalt.org/zap">move along</a>.
    </body>
</html>
`)
}

func TestPostRejected(t *testing.T) {
	t.Parallel()

	h := CreateHandler(&Config{
		URL: "go.uberalt.org",
		Packages: map[string]PackageConfig{
			"zap": {
				Repo: "github.com/uber-go/zap",
			},
		},
	})
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	tests := []struct {
		desc string
		path string
	}{
		{desc: "index", path: "/"},
		{desc: "package", path: "/zap"},
		{desc: "subpackage", path: "/zap/zapcore"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			res, err := http.Post(srv.URL+tt.path, "text/plain", strings.NewReader("foo"))
			require.NoError(t, err)
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusNotFound, res.StatusCode,
				"expected 404, got:\n%s", string(body))
		})
	}
}
