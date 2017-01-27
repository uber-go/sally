package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/sally"
)

var (
	yml  = flag.String("yml", "sally.yaml", "yaml file to read config from")
	port = flag.Int("port", 8080, "port to listen and serve on")
)

func main() {
	flag.Parse()

	log.Printf("Parsing yaml at path: %s\n", *yml)
	config, err := sally.Parse(*yml)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", *yml, err)
	}
	log.Printf("Creating HTTP handler with config: %v", config)
	handler := sally.CreateHandler(config)
	log.Printf(`Starting HTTP handler on ":%d"`, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
