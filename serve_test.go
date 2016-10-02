package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func GetHandlerFromYAML(t *testing.T, content string) http.Handler {
	var config Config
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		t.Fatal(err)
	}
	handler, err := GetHandler(config)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func Record(t *testing.T, config string, uri string) *httptest.ResponseRecorder {
	handler := GetHandlerFromYAML(t, config)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func TestNonExistentPackageShould404(t *testing.T) {
	rr := Record(t, "", "/nonexistent")
	assert.Equal(t, rr.Code, http.StatusNotFound)
}
