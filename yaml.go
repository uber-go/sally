package sally

import (
	"go.pedge.io/pkg/errors"
	"gopkg.in/yaml.v2"
)

// yamlConfig is the yaml configuation.
type yamlConfig struct {
	URL                  string                  `yaml:"url"`                     // the base url, not including https://
	Packages             map[string]*yamlPackage `yaml:"packages"`                // packages to redirect, map from path to package, the go package is URL/path
	IndexTitle           string                  `yaml:"index_title"`             // the title of the index page
	IndexLinks           []*yamlIndexLink        `yaml:"index_links"`             // extra links to display
	IndexShortcutIconURL string                  `yaml:"index_shortcut_icon_url"` // the shortcut icon urk for the index page
}

// yamlPackage is the yaml representation of a package.
type yamlPackage struct {
	Type       string   `yaml:"type"`        // the type, must be github
	GithubUser string   `yaml:"github_user"` // the github user, required if type is github
	GithubRepo string   `yaml:"github_repo"` // the github repository, required if type is github
	Badges     []string `yaml:"badges"`      // badges to display, must be from circleci, goreportcard, godoc, license-mit
	Private    bool     `yaml:"private"`     // if true, the page will not be displayed on the index page
}

// yamlIndexLink is the yaml representation of an extra link to display on the index page.
type yamlIndexLink struct {
	URL  string `yaml:"url"`  // the url to redirect to.
	Name string `yaml:"name"` // the name to display on the index page
}

// newYAMLConfig reads and validates a config from yaml data.
func newYAMLConfig(data []byte) (*yamlConfig, error) {
	yamlConfig := &yamlConfig{}
	if err := yaml.Unmarshal(data, yamlConfig); err != nil {
		return nil, err
	}
	if err := validateYAMLConfig(yamlConfig); err != nil {
		return nil, err
	}
	return yamlConfig, nil
}

func validateYAMLConfig(yamlConfig *yamlConfig) error {
	if yamlConfig.URL == "" {
		return pkgerrors.New("nil", "field_path", "url")
	}
	for _, yamlPackage := range yamlConfig.Packages {
		if err := validateYAMLPackage(yamlPackage); err != nil {
			return err
		}
	}
	return nil
}

func validateYAMLPackage(yamlPackage *yamlPackage) error {
	switch yamlPackage.Type {
	case "github":
		if yamlPackage.GithubUser == "" {
			return pkgerrors.New("nil", "field_path", "packages.github_user")
		}
		if yamlPackage.GithubRepo == "" {
			return pkgerrors.New("nil", "field_path", "packages.github_repo")
		}
	default:
		return pkgerrors.New("unknown package type", "type", yamlPackage.Type)
	}
	for _, yamlBadge := range yamlPackage.Badges {
		if err := validateYAMLBadge(yamlBadge); err != nil {
			return err
		}
	}
	return nil
}

func validateYAMLBadge(yamlBadge string) error {
	switch yamlBadge {
	case "circleci", "goreportcard", "godoc", "license-mit":
		return nil
	default:
		return pkgerrors.New("unknown badge type", "type", yamlBadge)
	}
}
