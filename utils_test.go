package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

// TempFile persists contents and returns the path and a clean func
func TempFile(t *testing.T, contents string) (path string, clean func()) {
	content := []byte(contents)
	tmpfile, err := os.CreateTemp("", "sally-tmp")
	if err != nil {
		t.Fatal("Unable to create tmpfile", err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal("Unable to write tmpfile", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Unable to close tmpfile", err)
	}

	return tmpfile.Name(), func() {
		_ = os.Remove(tmpfile.Name())
	}
}

// CreateHandlerFromYAML builds the Sally handler from a yaml config string
func CreateHandlerFromYAML(t *testing.T, content string) (handler http.Handler, clean func()) {
	path, clean := TempFile(t, content)

	config, err := Parse(path)
	if err != nil {
		t.Fatalf("Unable to parse %s: %v", path, err)
	}

	return CreateHandler(config), clean
}

// CallAndRecord makes a GET request to the Sally handler and returns a response recorder
func CallAndRecord(t *testing.T, config string, uri string) *httptest.ResponseRecorder {
	handler, clean := CreateHandlerFromYAML(t, config)
	defer clean()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatalf("Unable to create request to %s: %v", uri, err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

// AssertResponse normalizes and asserts the body from rr against want
func AssertResponse(t *testing.T, rr *httptest.ResponseRecorder, code int, want string) {
	assert.Equal(t, rr.Code, code)
	assert.Equal(t, reformatHTML(t, want), reformatHTML(t, rr.Body.String()))
}

func reformatHTML(t *testing.T, s string) string {
	n, err := html.Parse(strings.NewReader(s))
	require.NoError(t, err)

	var buff bytes.Buffer
	require.NoError(t, html.Render(&buff, n))
	return buff.String()
}
