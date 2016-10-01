package main

import (
	"flag"
	"log"
)

func main() {
	yml := flag.String("yml", "sally.yaml", "yaml file to read config from")
	flag.Parse()

	config, err := Parse(*yml)
	if err != nil {
		log.Fatal(err)
	}

	if err := Serve(config); err != nil {
		log.Fatal(err)
	}

	select {}
}
