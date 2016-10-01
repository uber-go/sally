package main

import (
	"fmt"
	"html"
	"net/http"
)

// Serve starts the HTTP server
func Serve(config Config) error {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}
