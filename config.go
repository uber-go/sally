package main

import (
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

const (
	_defaultGodocServer = "pkg.go.dev"
)

// Config defines the configuration for a Sally server.
type Config struct {
	// URL is the base URL for all vanity imports.
	URL string `yaml:"url"` // required

	// Packages is a map of package name to package details.
	Packages map[string]PackageConfig `yaml:"packages"`

	// Godoc specifies where to redirect to for documentation.
	Godoc GodocConfig `yaml:"godoc"`
}

// GodocConfig is the configuration for the documentation server.
type GodocConfig struct {
	// Host is the hostname of the documentation server.
	//
	// Defaults to pkg.go.dev.
	Host string `yaml:"host"`
}

// PackageConfig is the configuration for a single Go module
// that is served by Sally.
type PackageConfig struct {
	// Repo is the URL to the Git repository for the module
	// without the https:// prefix.
	// This URL must serve the Git HTTPS protocol.
	//
	// For example, "github.com/uber-go/sally".
	Repo string `yaml:"repo"` // required

	// URL is the base URL of the vanity import for this module.
	//
	// Defaults to the URL specified in the top-level config.
	URL string `yaml:"url"`

	// Desc is a plain text description of this module.
	Desc string `yaml:"description"`
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

	return &c, err
}
