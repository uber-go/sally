package main

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

// TempFile persists contents and returns the path and a clean func
func TempFile(t *testing.T, contents string) (path string) {
	content := []byte(contents)
	tmpfile, err := os.CreateTemp("", "sally-tmp")
	require.NoError(t, err, "unable to create tmpfile")

	_, err = tmpfile.Write(content)
	require.NoError(t, err, "unable to write tmpfile")

	err = tmpfile.Close()
	require.NoError(t, err, "unable to close tmpfile")

	t.Cleanup(func() {
		_ = os.Remove(tmpfile.Name())
	})

	return tmpfile.Name()
}

// CreateHandlerFromYAML builds the Sally handler from a yaml config string
func CreateHandlerFromYAML(t *testing.T, templates *template.Template, content string) (handler http.Handler) {
	path := TempFile(t, content)

	config, err := Parse(path)
	require.NoError(t, err, "unable to parse path %s", path)

	handler, err = CreateHandler(config, templates)
	require.NoError(t, err)

	return handler
}

// CallAndRecord makes a GET request to the Sally handler and returns a response recorder
func CallAndRecord(t *testing.T, config string, templates *template.Template, uri string) *httptest.ResponseRecorder {
	handler := CreateHandlerFromYAML(t, templates, config)

	req, err := http.NewRequest("GET", uri, nil)
	require.NoError(t, err, "unable to create request to %s", uri)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

// AssertResponse normalizes and asserts the body from rr against want
func AssertResponse(t *testing.T, rr *httptest.ResponseRecorder, code int, want string) {
	assert.Equal(t, rr.Code, code)
	assert.Equal(t, reformatHTML(t, want), reformatHTML(t, rr.Body.String()))
}

// getTestTemplates returns a [template.Template] object with the default templates,
// overwritten by the  [overrideTemplates]. If [overrideTemplates] is nil, the returned
// templates are a clone of the global [_templates].
func getTestTemplates(tb testing.TB, overrideTemplates map[string]string) *template.Template {
	if len(overrideTemplates) == 0 {
		// We must clone! Cloning can only be done before templates are executed. Therefore,
		// we cannot run some tests without cloning, and then attempt cloning it. It'll panic.
		templates, err := _templates.Clone()
		require.NoError(tb, err)
		return templates
	}

	templatesDir := tb.TempDir() // This is automatically removed at the end of the test.
	for name, content := range overrideTemplates {
		err := os.WriteFile(filepath.Join(templatesDir, name), []byte(content), 0o666)
		require.NoError(tb, err)
	}

	templates, err := getCombinedTemplates(templatesDir)
	require.NoError(tb, err)
	return templates
}

func reformatHTML(t *testing.T, s string) string {
	n, err := html.Parse(strings.NewReader(s))
	require.NoError(t, err)

	var buff bytes.Buffer
	require.NoError(t, html.Render(&buff, n))
	return buff.String()
}
