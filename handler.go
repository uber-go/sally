package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var indexTemplate, packageTemplate *template.Template

func init() {
	tmpls := template.Must(template.ParseGlob("templates/*.html"))
	if indexTemplate = tmpls.Lookup("index.html"); indexTemplate == nil {
		panic("Missing index.html template")
	}
	if packageTemplate = tmpls.Lookup("package.html"); packageTemplate == nil {
		panic("Missing package.html template")
	}
}

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
	canonicalURL := fmt.Sprintf("%s/%s", h.config.URL, h.pkgName)
	data := struct {
		Repo         string
		CanonicalURL string
		GodocURL     string
	}{
		Repo:         h.pkg.Repo,
		CanonicalURL: canonicalURL,
		GodocURL:     fmt.Sprintf("https://godoc.org/%s%s", canonicalURL, ps.ByName("path")),
	}
	if err := packageTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
