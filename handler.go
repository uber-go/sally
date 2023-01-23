package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"go.uber.org/sally/templates"
)

var (
	indexTemplate = template.Must(
		template.New("index.html").Parse(templates.Index))
	packageTemplate = template.Must(
		template.New("package.html").Parse(templates.Package))
)

// CreateHandler creates a Sally http.Handler
func CreateHandler(config *Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", &indexHandler{config: config})
	for name, pkg := range config.Packages {
		handle := packageHandler{
			pkgName: name,
			pkg:     pkg,
			config:  config,
		}
		// Double-register so that "/foo"
		// does not redirect to "/foo/" with a 300.
		mux.Handle("/"+name, &handle)
		mux.Handle("/"+name+"/", &handle)
	}

	return mux
}

type indexHandler struct {
	config *Config
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Index handler only supports '/'.
	// ServeMux will call us for any '/foo' that is not a known package.
	if r.Method != http.MethodGet || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := indexTemplate.Execute(w, h.config); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type packageHandler struct {
	pkgName string
	pkg     Package
	config  *Config
}

func (h *packageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Extract the relative path to subpackages, if any.
	//	"/foo/bar" => "/bar"
	//	"/foo" => ""
	relPath := strings.TrimPrefix(r.URL.Path, "/"+h.pkgName)

	baseURL := h.config.URL
	if h.pkg.URL != "" {
		baseURL = h.pkg.URL
	}
	canonicalURL := fmt.Sprintf("%s/%s", baseURL, h.pkgName)
	data := struct {
		Repo         string
		Branch       string
		CanonicalURL string
		GodocURL     string
	}{
		Repo:         h.pkg.Repo,
		Branch:       h.pkg.Branch,
		CanonicalURL: canonicalURL,
		GodocURL:     fmt.Sprintf("https://%s/%s%s", h.config.Godoc.Host, canonicalURL, relPath),
	}
	if err := packageTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
