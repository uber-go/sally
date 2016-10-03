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
	flag.Parse()

	log.Printf("Parsing yaml at path: %s\n", *yml)
	config, err := Parse(*yml)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating HTTP handler with config: %s", config)
	handler := CreateHandler(config)

	log.Printf(`Starting HTTP handler on ":%d"`, *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler); err != nil {
		log.Fatal(err)
	}
}
