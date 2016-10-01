package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ListenAndServe starts a Sally server
func ListenAndServe(port int, config Config) error {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	handle, err := index(config)
	if err != nil {
		return err
	}
	router.GET("/", handle)

	for name, p := range config.Packages {
		handle, err := pkg(pkgViewModel{
			Package: p,
			Name:    name,
			Config:  config,
		})
		if err != nil {
			return err
		}
		router.GET(fmt.Sprintf("/%s", name), handle)
		router.GET(fmt.Sprintf("/%s/*name", name), handle)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		return err
	}

	return nil
}

func index(config Config) (httprouter.Handle, error) {
	t, err := template.New("index").Parse(`
<h1>Hello World</h1>
`)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t.Execute(w, config)
	}, nil
}

func pkg(p pkgViewModel) (httprouter.Handle, error) {
	t, err := template.New("package").Parse(`
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
`)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t.Execute(w, p)
	}, nil
}

type pkgViewModel struct {
	Package

	Name   string
	Config Config
}

func (p pkgViewModel) CanonicalURL() string {
	return fmt.Sprintf("%s/%s", p.Config.URL, p.Name)
}

func (p pkgViewModel) GodocURL() string {
	return fmt.Sprintf("https://godoc.org/%s", p.CanonicalURL())
}
