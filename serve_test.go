package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetHandlerFromYAML(t *testing.T, content string) http.Handler {
	// TODO pass in yaml from tests by using ioutil.TempFile()
	config, err := Parse("sally.yaml")
	if err != nil {
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

var config = `

	url: go.uber.org
	packages:
	  yarpc:
	    repo: github.com/yarpc/yarpc-go

`

func TestPackageShouldExist(t *testing.T) {
	rr := Record(t, config, "/yarpc")
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
	rr := Record(t, config, "/nonexistent")
	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestTrailingSlash(t *testing.T) {
	rr := Record(t, config, "/yarpc/")
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
	rr := Record(t, config, "/yarpc/heeheehee")
	assert.Equal(t, rr.Code, http.StatusOK)

	rr = Record(t, config, "/yarpc/heehee/hawhaw")
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
