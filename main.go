package main

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

	config, err := Parse(*yml)
	if err != nil {
		log.Fatal(err)
	}

	handler, err := GetHandler(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler); err != nil {
		log.Fatal(err)
	}
}
