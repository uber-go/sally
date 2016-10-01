package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Serve starts the HTTP server
func Serve(config Config) error {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	router.GET("/", index(config))

	for name, p := range config.Packages {
		h := pkg(p)
		router.GET(fmt.Sprintf("/%s", name), h)
		router.GET(fmt.Sprintf("/%s/*name", name), h)
	}

	// TODO port should be cli opt
	if err := http.ListenAndServe(":8080", router); err != nil {
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
