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
		handle, err := pkg(p)
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

func pkg(pkg Package) (httprouter.Handle, error) {
	t, err := template.New("package").Parse(`
<h1>Package</h1>
`)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t.Execute(w, pkg)
	}, nil
}
