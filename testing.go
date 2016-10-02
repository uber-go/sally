package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TempFile persists contents and returns the path and a clean func
func TempFile(t *testing.T, contents string) (path string, clean func()) {
	content := []byte(contents)
	tmpfile, err := ioutil.TempFile("", "sally-tmp")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name(), func() {
		os.Remove(tmpfile.Name())
	}
}

// GetHandlerFromYAML builds the Sally handler from a yaml config string
func GetHandlerFromYAML(t *testing.T, content string) (handler http.Handler, clean func()) {
	path, clean := TempFile(t, content)

	config, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	handler, err = GetHandler(config)
	if err != nil {
		t.Fatal(err)
	}

	return handler, clean
}

// Record makes a GET request to the Sally handler and returns a response recorder
func Record(t *testing.T, config string, uri string) *httptest.ResponseRecorder {
	handler, clean := GetHandlerFromYAML(t, config)
	defer clean()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}
