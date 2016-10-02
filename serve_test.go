package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func WithConfig(content string, fn func(handler http.Handler)) {
	fmt.Println(content)
	var config Config
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		log.Panic(err)
	}
	handler, err := GetHandler(config)
	if err != nil {
		log.Panic(err)
	}
	fn(handler)
}

func TestNonExistentPackageShould404(t *testing.T) {
	WithConfig("", func(handler http.Handler) {
		req, err := http.NewRequest("GET", "/nonexistent", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusNotFound)
	})
}
