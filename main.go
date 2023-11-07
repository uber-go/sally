// sally is an HTTP service that allows you to serve
// vanity import paths for your Go packages.
package main // import "go.uber.org/sally"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	yml := flag.String("yml", "sally.yaml", "yaml file to read config from")
	tpls := flag.String("templates", "", "directory of .html templates to use")
	port := flag.Int("port", 8080, "port to listen and serve on")
	flag.Parse()

	log.Printf("Parsing yaml at path: %s\n", *yml)
	config, err := Parse(*yml)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", *yml, err)
	}

	if *tpls != "" {
		log.Printf("Parsing templates at path: %s\n", *tpls)
		templates, err = templates.ParseGlob(filepath.Join(*tpls, "*.html"))
		if err != nil {
			log.Fatalf("Failed to parse templates at %s: %v", *tpls, err)
		}
	}

	log.Printf("Creating HTTP handler with config: %v", config)
	handler := CreateHandler(config)

	log.Printf(`Starting HTTP handler on ":%d"`, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
