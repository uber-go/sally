package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	router.GET("/", indexHandler{config: config}.Handle)

	for name, pkg := range config.Packages {
		handle := packageHandler{
			pkgName: name,
			pkg:     pkg,
			config:  config,
		}.Handle
		router.GET(fmt.Sprintf("/%s", name), handle)
		router.GET(fmt.Sprintf("/%s/*path", name), handle)
	}

	return router
}

type indexHandler struct {
	config *Config
}

func (h indexHandler) Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := indexTemplate.Execute(w, h.config); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type packageHandler struct {
	pkgName string
	pkg     Package
	config  *Config
}

func (h packageHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		GodocURL:     fmt.Sprintf("https://%s/%s%s", h.config.Godoc.Host, canonicalURL, ps.ByName("path")),
	}
	if err := packageTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
