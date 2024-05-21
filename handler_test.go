package main

import (
	"html/template"
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
  scago:
    repo: github.com/m5ka/scago
    doc_url: https://example.org/docs/go-pkg/scago
    doc_badge: https://img.shields.io/badge/custom_docs-scago-blue?logo=go

`

func TestIndex(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/")
	assert.Equal(t, 200, rr.Code)

	body := rr.Body.String()
	assert.Contains(t, body, "github.com/thriftrw/thriftrw-go")
	assert.Contains(t, body, "github.com/yarpc/yarpc-go")
	assert.Contains(t, body, "A fast, structured logging library.")
	assert.Contains(t, body, "github.com/yarpc/metrics")
	assert.Contains(t, body, "github.com/yarpc/something")
	assert.Contains(t, body, "github.com/m5ka/scago")
}

func TestSubindex(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/net")
	assert.Equal(t, 200, rr.Code)

	body := rr.Body.String()
	assert.NotContains(t, body, "github.com/thriftrw/thriftrw-go")
	assert.NotContains(t, body, "github.com/m5ka/scago")
	assert.NotContains(t, body, "github.com/yarpc/yarpc-go")
	assert.Contains(t, body, "github.com/yarpc/metrics")
	assert.Contains(t, body, "github.com/yarpc/something")
}

func TestPackageShouldExist(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/yarpc")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc">
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
                a { color: #ddd; }
                a:visited { color: #bbb; }
            }
        </style>
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc">move along</a>.
    </body>
</html>
`)
}

func TestNonExistentPackageShould404(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/nonexistent")
	assert.Equal(t, "no-cache", rr.Header().Get("Cache-Control"))
	AssertResponse(t, rr, 404, `<!DOCTYPE html>
<html>
    <head>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" />
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
            }
        </style>
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
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/yarpc/")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc/">
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
                a { color: #ddd; }
                a:visited { color: #bbb; }
            }
        </style>
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc/">move along</a>.
    </body>
</html>
`)
}

func TestDeepImports(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/yarpc/heeheehee")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc/heeheehee">
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
                a { color: #ddd; }
                a:visited { color: #bbb; }
            }
        </style>
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc/heeheehee">move along</a>.
    </body>
</html>
`)

	rr = CallAndRecord(t, config, getTestTemplates(t, nil), "/yarpc/heehee/hawhaw")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/yarpc git https://github.com/yarpc/yarpc-go">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uber.org/yarpc/heehee/hawhaw">
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
                a { color: #ddd; }
                a:visited { color: #bbb; }
            }
        </style>
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uber.org/yarpc/heehee/hawhaw">move along</a>.
    </body>
</html>
`)
}

func TestPackageLevelURL(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/zap")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uberalt.org/zap git https://github.com/uber-go/zap">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uberalt.org/zap">
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
                a { color: #ddd; }
                a:visited { color: #bbb; }
            }
        </style>
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uberalt.org/zap">move along</a>.
    </body>
</html>
`)
}

func TestCustomDocURL(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/scago")
	AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uber.org/scago git https://github.com/m5ka/scago">
        <meta http-equiv="refresh" content="0; url=https://example.org/docs/go-pkg/scago">
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
                a { color: #ddd; }
                a:visited { color: #bbb; }
            }
        </style>
    </head>
    <body>
        Nothing to see here. Please <a href="https://example.org/docs/go-pkg/scago">move along</a>.
    </body>
</html>
`)
}

func TestCustomDocBadge(t *testing.T) {
	rr := CallAndRecord(t, config, getTestTemplates(t, nil), "/")
	assert.Equal(t, 200, rr.Code)

	body := rr.Body.String()
	assert.Contains(t, body, "<img src=\"//pkg.go.dev/badge/go.uber.org/yarpc.svg\" alt=\"Go Reference\" />")
	assert.Contains(t, body, "<img src=\"//pkg.go.dev/badge/go.uberalt.org/zap.svg\" alt=\"Go Reference\" />")
	assert.Contains(t, body,
		"<img src=\"https://img.shields.io/badge/custom_docs-scago-blue?logo=go\" alt=\"Go Reference\" />")
	assert.NotContains(t, body, "<img src=\"//pkg.go.dev/badge/go.uber.org/scago.svg\" alt=\"Go Reference\" />")
}

func TestPostRejected(t *testing.T) {
	t.Parallel()

	h, err := CreateHandler(&Config{
		URL: "go.uberalt.org",
		Packages: map[string]PackageConfig{
			"zap": {
				Repo: "github.com/uber-go/zap",
			},
		},
	}, getTestTemplates(t, nil))
	require.NoError(t, err)
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
			templates := getTestTemplates(t, nil)
			h := newIndexHandler(tt.pkgs, templates.Lookup("index.html"), templates.Lookup("404.html"))
			start, end := h.rangeOf(tt.path)

			var got []string
			for _, pkg := range tt.pkgs[start:end] {
				got = append(got, pkg.Name)
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestCustomTemplates(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		for _, name := range []string{"index.html", "package.html", "404.html"} {
			templatesText := map[string]string{
				"index.html":   "index",
				"package.html": "package",
				"404.html":     "404",
			}
			delete(templatesText, name)

			templates := template.New("")
			for tplName, tplText := range templatesText {
				var err error
				templates, err = templates.New(tplName).Parse(tplText)
				require.NoError(t, err)
			}

			_, err := CreateHandler(&Config{}, templates)
			require.Error(t, err, name)
		}
	})

	t.Run("replace", func(t *testing.T) {
		templates := getTestTemplates(t, map[string]string{
			"404.html": "not found: {{ .Path }}",
		})

		// Overrides 404.html
		rr := CallAndRecord(t, config, templates, "/blah")
		require.Equal(t, http.StatusNotFound, rr.Result().StatusCode)

		// But not package.html
		rr = CallAndRecord(t, config, templates, "/zap")
		AssertResponse(t, rr, 200, `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="go.uberalt.org/zap git https://github.com/uber-go/zap">
        <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/go.uberalt.org/zap">
        <style>
            @media (prefers-color-scheme: dark) {
                body { background-color: #333; color: #ddd; }
                a { color: #ddd; }
                a:visited { color: #bbb; }
            }
        </style>
    </head>
    <body>
        Nothing to see here. Please <a href="https://pkg.go.dev/go.uberalt.org/zap">move along</a>.
    </body>
</html>
`)
	})
}

func BenchmarkHandlerDispatch(b *testing.B) {
	handler, err := CreateHandler(&Config{
		URL: "go.uberalt.org",
		Packages: map[string]PackageConfig{
			"zap": {
				Repo: "github.com/uber-go/zap",
			},
			"net/metrics": {
				Repo: "github.com/yarpc/metrics",
			},
		},
	}, getTestTemplates(b, nil))
	require.NoError(b, err)
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

func (nopResponseWriter) Header() http.Header       { return http.Header{} }
func (nopResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (nopResponseWriter) WriteHeader(int)           {}
