package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"go.uber.org/sally/templates"
)

var (
	indexTemplate = template.Must(
		template.New("index.html").Parse(templates.Index))
	packageTemplate = template.Must(
		template.New("package.html").Parse(templates.Package))
)

// Handler handles inbound HTTP requests.
//
// It provides the following endpoints:
//
//	GET /
//		Index page listing all packages.
//	GET /<name>
//		Package page for the given package.
//	GET /<dir>
//		Page listing packages under the given directory,
//		assuming that there's no package with the given name.
//	GET /<name>/<subpkg>
//		Package page for the given subpackage.
type Handler struct {
	pkgs pathTree[*sallyPackage]
}

// CreateHandler builds a new handler
// with the provided package configuration.
func CreateHandler(config *Config) *Handler {
	var pkgs pathTree[*sallyPackage]
	for name, pkg := range config.Packages {
		baseURL := config.URL
		if pkg.URL != "" {
			// Package-specific override for the base URL.
			baseURL = pkg.URL
		}
		modulePath := path.Join(baseURL, name)
		docURL := "https://" + path.Join(config.Godoc.Host, modulePath)

		pkgs.Set(name, &sallyPackage{
			Desc:       pkg.Desc,
			ModulePath: modulePath,
			DocURL:     docURL,
			GitURL:     pkg.Repo,
		})
	}

	return &Handler{
		pkgs: pkgs,
	}
}

var _ http.Handler = (*Handler)(nil)

type sallyPackage struct {
	// Canonical import path for the package.
	ModulePath string

	// Description of the package, if any.
	Desc string

	// URL at which documentation for the package can be found.
	DocURL string

	// URL at which the Git repository is hosted.
	GitURL string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.serveHTTP(w, r); err != nil {
		if herr := new(httpError); errors.As(err, &herr) {
			http.Error(w, herr.Message, herr.Code)
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
}

// httpError indicates that an HTTP error occurred.
//
// The caller will write the error code and message to the response.
type httpError struct {
	Code    int    // HTTP status code
	Message string // error message
}

func httpErrorf(code int, format string, args ...interface{}) error {
	return &httpError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func (e *httpError) Error() string {
	return fmt.Sprintf("status %d: %s", e.Code, e.Message)
}

// serveHTTP is similar to ServeHTTP, except it returns an error.
//
// If it returns an httpError,
// the caller will write the error code and message to the response.
func (h *Handler) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httpErrorf(http.StatusNotFound, "method %q not allowed", r.Method)
	}

	path := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/")

	if pkg, suffix, ok := h.pkgs.Lookup(path); ok {
		return h.servePackage(w, pkg, suffix)
	}
	return h.serveIndex(w, path, h.pkgs.ListByPath(path))
}

func (h *Handler) servePackage(w http.ResponseWriter, pkg *sallyPackage, suffix string) error {
	return packageTemplate.Execute(w,
		struct {
			ModulePath string
			GitURL     string
			DocURL     string
		}{
			ModulePath: pkg.ModulePath,
			GitURL:     pkg.GitURL,
			DocURL:     pkg.DocURL + suffix,
		})
}

func (h *Handler) serveIndex(w http.ResponseWriter, path string, pkgs []*sallyPackage) error {
	if len(pkgs) == 0 {
		return httpErrorf(http.StatusNotFound, "no packages found under path: %s", path)
	}

	return indexTemplate.Execute(w,
		struct{ Packages []*sallyPackage }{
			Packages: pkgs,
		})
}
