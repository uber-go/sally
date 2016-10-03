package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var indexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
    <body>
        <ul>
            {{ range $key, $value := .Packages }}
	  	        <li>{{ $key }} - {{ $value.Repo }}</li>
	        {{ end }}
        </ul>
    </body>
</html>
`))

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

// CreateHandler creates a Sally http.Handler
func CreateHandler(config Config) http.Handler {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	router.GET("/", indexHandler{ViewModel: config}.Handle)

	for name, pkg := range config.Packages {
		handle := packageHandler{
			ViewModel: packageViewModel{
				Package: pkg,
				Name:    name,
				Config:  config,
			},
		}.Handle
		router.GET(fmt.Sprintf("/%s", name), handle)
		router.GET(fmt.Sprintf("/%s/*path", name), handle)
	}

	return router
}

type indexHandler struct {
	ViewModel Config
}

func (h indexHandler) Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	indexTemplate.Execute(w, h.ViewModel)
}

type packageHandler struct {
	ViewModel packageViewModel
}

func (h packageHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	packageTemplate.Execute(w, h.ViewModel.NewWithAddlGodocPath(ps.ByName("path")))
}

type packageViewModel struct {
	Package

	Name          string
	Config        Config
	AddlGodocPath string
}

func (p packageViewModel) CanonicalURL() string {
	return fmt.Sprintf("%s/%s", p.Config.URL, p.Name)
}

func (p packageViewModel) GodocURL() string {
	return fmt.Sprintf("https://godoc.org/%s%s", p.CanonicalURL(), p.AddlGodocPath)
}

func (p packageViewModel) NewWithAddlGodocPath(uri string) packageViewModel {
	p.AddlGodocPath = uri
	return p
}
