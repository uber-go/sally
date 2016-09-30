package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config represents the structure of the yaml file
type Config struct {
	URL      string             `yaml:"url"`
	Packages map[string]Package `yaml:"packages"`
}

// Package details the options available for each repo
type Package struct {
	Repo    string `yaml:"repo"`
	Private bool   `yaml:"private"`
}

var _privateGHReplacer = strings.NewReplacer("github.com/", "")

// FetchURL returns the correct fetch URL, respecting if a repo is public/private
func (p Package) FetchURL() string {
	if p.Private {
		return p.PrivateURL()
	}
	return p.PublicURL()
}

// PrivateURL returns the URL you'd use to pull this repo if it's still private (e.g. in development)
func (p Package) PrivateURL() string {
	// TODO(ai) currently only works for GitHub
	return fmt.Sprintf("git@github.com:%s", _privateGHReplacer.Replace(p.Repo))
}

// PublicURL returns a URL to fetch the source of a package
func (p Package) PublicURL() string {
	// TODO(ai) currently only works for GitHub
	return "https://" + p.Repo
}

// Parse takes a path to a yaml file and produces a parsed Config
func Parse(path string) (Config, error) {
	var c Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, err
	}

	return c, err
}
