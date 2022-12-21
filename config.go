package main

import (
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	_defaultGodocServer = "pkg.go.dev"
	_defaultBranch      = "master"
)

// Config represents the structure of the yaml file
type Config struct {
	URL      string             `yaml:"url"`
	Packages map[string]Package `yaml:"packages"`
	Godoc    struct {
		Host string `yaml:"host"`
	} `yaml:"godoc"`
}

// Package details the options available for each repo
type Package struct {
	Repo   string `yaml:"repo"`
	Branch string `yaml:"branch"`
	URL    string `yaml:"url"`
}

// Parse takes a path to a yaml file and produces a parsed Config
func Parse(path string) (*Config, error) {
	var c Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	if c.Godoc.Host == "" {
		c.Godoc.Host = _defaultGodocServer
	} else {
		host := c.Godoc.Host
		host = strings.TrimPrefix(host, "https://")
		host = strings.TrimPrefix(host, "http://")
		host = strings.TrimSuffix(host, "/")
		c.Godoc.Host = host
	}

	// set default branch
	for v, p := range c.Packages {
		if p.Branch == "" {
			p.Branch = _defaultBranch
			c.Packages[v] = p
		}
	}

	return &c, err
}
