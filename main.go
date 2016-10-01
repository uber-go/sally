package main

import (
	"flag"
	"log"
)

// Config represents the structure of the yaml file
type Config struct {
	URL      string             `yaml:"url"`
	Packages map[string]Package `yaml:"packages"`
}

// Package details the options available for each repo
type Package struct {
	Repo string `yaml:"repo"`
}

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
