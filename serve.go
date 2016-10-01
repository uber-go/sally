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

	router.GET("/", handleIndex)

	for name, pkg := range config.Packages {
		fmt.Println(pkg)
		router.GET(fmt.Sprintf("/%s", name), handlePackage)
		router.GET(fmt.Sprintf("/%s/*name", name), handlePackage)
	}

	// TODO port should be cli opt
	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}

func handleIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func handlePackage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
