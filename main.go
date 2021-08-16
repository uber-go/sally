package main // import "go.uber.org/sally"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	yml := flag.String("yml", "sally.yaml", "yaml file to read config from")
	port := flag.Int("port", 8080, "port to listen and serve on")
	output := flag.String("o", "", "generate static site to directory")
	flag.Parse()

	log.Printf("Parsing yaml at path: %s\n", *yml)
	config, err := Parse(*yml)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", *yml, err)
	}

	if output != nil && *output != "" {
		err = GenerateSite(config, *output)
		if err != nil {
			log.Fatalf("Failed to generate static site: %v", err)
		}
		return
	}

	log.Printf("Creating HTTP handler with config: %v", config)
	handler := CreateHandler(config)

	log.Printf(`Starting HTTP handler on ":%d"`, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
