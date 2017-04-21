package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"
)

// CreateHandler creates a Sally http.Handler
func CreateHandler(config *Config) http.Handler {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	router.GET("/", indexHandler{config: newIndexConfig(config)}.Handle)

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
	config *indexConfig
}

func (h indexHandler) Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := indexTemplate.Execute(w, h.config); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

var indexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
    <body>
        <ul>
            {{ range $package := .Packages }}
	  	        <li>{{ $package.Name }} - {{ $package.Repo }}</li>
	        {{ end }}
        </ul>
    </body>
</html>
`))

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

var packageTemplate = template.Must(template.New("package").Parse(`
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="{{ .CanonicalURL }} git https://{{ .Repo }}">
        <meta name="go-source" content="{{ .CanonicalURL }} https://{{ .Repo }} https://{{ .Repo }}/tree/master{/dir} https://{{ .Repo }}/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url={{ .GodocURL }}">
    </head>
    <body>
        Nothing to see here. Please <a href="{{ .GodocURL }}">move along</a>.
    </body>
</html>
`))

type indexConfig struct {
	Packages []*indexPackage
}

type indexPackage struct {
	Name string
	Repo string
}

func newIndexConfig(c *Config) *indexConfig {
	i := &indexConfig{}
	for name, pkg := range c.Packages {
		i.Packages = append(i.Packages, &indexPackage{name, pkg.Repo})
	}
	sort.Sort(indexPackagesByName(i.Packages))
	return i
}

type indexPackagesByName []*indexPackage

func (p indexPackagesByName) Len() int           { return len(p) }
func (p indexPackagesByName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p indexPackagesByName) Less(i, j int) bool { return p[i].Name < p[j].Name }
