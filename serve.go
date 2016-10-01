package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ListenAndServe starts a Sally server
func ListenAndServe(port int, config Config) error {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	router.GET("/", index(config))

	for name, p := range config.Packages {
		h := pkg(p)
		router.GET(fmt.Sprintf("/%s", name), h)
		router.GET(fmt.Sprintf("/%s/*name", name), h)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		return err
	}

	return nil
}

func index(config Config) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Welcome!\n")
	}
}

func pkg(pkg Package) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	}
}
