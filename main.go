package main

import (
	"flag"
	"log"
)

func main() {
	yml := flag.String("yml", "sally.yaml", "yaml file to read config from")
	port := flag.Int("port", 8080, "port to listen and serve on")
	flag.Parse()

	config, err := Parse(*yml)
	if err != nil {
		log.Fatal(err)
	}

	if err := ListenAndServe(*port, config); err != nil {
		log.Fatal(err)
	}

	select {}
}
