package main // import "go.uber.org/sally"

import (
	"flag"
	"log"
)

//go:generate go-bindata templates/

func main() {
	yml := flag.String("yml", "sally.yaml", "yaml file to read config from")
	dir := flag.String("dir", "out", "directory to write html files to")
	flag.Parse()

	c, err := Parse(*yml)
	if err != nil {
		log.Fatal(err)
	}

	if err := Write(c, *dir); err != nil {
		log.Fatal(err)
	}
}
