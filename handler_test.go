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
  net/metrics:
    repo: github.com/yarpc/metrics
  net/something:
    repo: github.com/yarpc/something

`

func TestIndex(t *testing.T) {
	rr := CallAndRecord(t, config, "/")
	assert.Equal(t, 200, rr.Code)

	body := rr.Body.String()
	assert.Contains(t, body, "github.com/thriftrw/thriftrw-go")
	assert.Contains(t, body, "github.com/yarpc/yarpc-go")
	assert.Contains(t, body, "A fast, structured logging library.")
	assert.Contains(t, body, "github.com/yarpc/metrics")
	assert.Contains(t, body, "github.com/yarpc/something")
}

func TestSubindex(t *testing.T) {
	rr := CallAndRecord(t, config, "/net")
	assert.Equal(t, 200, rr.Code)

	body := rr.Body.String()
	assert.NotContains(t, body, "github.com/thriftrw/thriftrw-go")
	assert.NotContains(t, body, "github.com/yarpc/yarpc-go")
	assert.Contains(t, body, "github.com/yarpc/metrics")
	assert.Contains(t, body, "github.com/yarpc/something")
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
	assert.Equal(t, "no-cache", rr.Header().Get("Cache-Control"))
	AssertResponse(t, rr, 404, `<!DOCTYPE html>
<html>
    <head>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" />
    </head>
    <body>
        <div class="container">
            <p>No packages found under: "nonexistent".</p>
        </div>
    </body>
</html>
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
			defer func() {
				assert.NoError(t, res.Body.Close())
			}()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusNotFound, res.StatusCode,
				"expected 404, got:\n%s", string(body))
		})
	}
}

func TestIndexHandler_rangeOf(t *testing.T) {
	tests := []struct {
		desc string
		pkgs []*sallyPackage
		path string
		want []string // names
	}{
		{
			desc: "empty",
			pkgs: []*sallyPackage{
				{Name: "foo"},
				{Name: "bar"},
			},
			want: []string{"foo", "bar"},
		},
		{
			desc: "single child",
			pkgs: []*sallyPackage{
				{Name: "foo/bar"},
				{Name: "baz"},
			},
			path: "foo",
			want: []string{"foo/bar"},
		},
		{
			desc: "multiple children",
			pkgs: []*sallyPackage{
				{Name: "foo/bar"},
				{Name: "foo/baz"},
				{Name: "qux"},
				{Name: "quux/quuz"},
			},
			path: "foo",
			want: []string{"foo/bar", "foo/baz"},
		},
		{
			desc: "to end of list",
			pkgs: []*sallyPackage{
				{Name: "a"},
				{Name: "b"},
				{Name: "c/d"},
				{Name: "c/e"},
			},
			path: "c",
			want: []string{"c/d", "c/e"},
		},
		{
			desc: "similar name",
			pkgs: []*sallyPackage{
				{Name: "foobar"},
				{Name: "foo/bar"},
			},
			path: "foo",
			want: []string{"foo/bar"},
		},
		{
			desc: "no match",
			pkgs: []*sallyPackage{
				{Name: "foo"},
				{Name: "bar"},
			},
			path: "baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			h := newIndexHandler(tt.pkgs)
			start, end := h.rangeOf(tt.path)

			var got []string
			for _, pkg := range tt.pkgs[start:end] {
				got = append(got, pkg.Name)
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func BenchmarkHandlerDispatch(b *testing.B) {
	handler := CreateHandler(&Config{
		URL: "go.uberalt.org",
		Packages: map[string]PackageConfig{
			"zap": {
				Repo: "github.com/uber-go/zap",
			},
			"net/metrics": {
				Repo: "github.com/yarpc/metrics",
			},
		},
	})
	resw := new(nopResponseWriter)

	tests := []struct {
		name string
		path string
	}{
		{name: "index", path: "/"},
		{name: "subindex", path: "/net"},
		{name: "package", path: "/zap"},
		{name: "subpackage", path: "/zap/zapcore"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			req := httptest.NewRequest("GET", tt.path, nil)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				handler.ServeHTTP(resw, req)
			}
		})
	}
}

type nopResponseWriter struct{}

func (nopResponseWriter) Header() http.Header       { return nil }
func (nopResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (nopResponseWriter) WriteHeader(int)           {}
