package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the structure of the yaml file
type Config struct {
	URL      string
	Packages map[string]Package
}

// Package details the options available for each repo
type Package struct {
	Repo string
}

// Parse takes a path to a yaml file and produces a parsed Config
func Parse(path string) (Config, error) {
	c := Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, err
}
