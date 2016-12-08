package sally

import (
	"fmt"

	"go.pedge.io/pkg/errors"
)

// tmplConfig is the configuration in a template-friendly manner.
type tmplConfig struct {
	Packages             []*tmplPackage
	IndexTitle           string
	IndexLinks           []*tmplIndexLink
	IndexShortcutIconURL string
}

// tmplPackage is a package to dispay.
type tmplPackage struct {
	Path       string       // the path to redirect from
	GoPackage  string       // the go package, basically ConfigURL/Path
	IndexURL   string       // the url to direct to on the index page
	VCS        string       // the vcs
	Repository string       // the repoitory to direct to
	Images     []*tmplImage // extra images to display on the index page
	Private    bool         // if true, page will not be displayed on the index page
}

// tmplImage is an image to display next to a package.
type tmplImage struct {
	URL    string // the url to redirect to
	Source string // the source url of the image
}

// tmplIndexLink is a template-friendly index link.
type tmplIndexLink struct {
	URL  string
	Name string
}

// newTmplConfig creates a new tmplConfig from a valid yamlConfig.
func newTmplConfig(yamlConfig *yamlConfig) (*tmplConfig, error) {
	tmplConfig := &tmplConfig{
		IndexTitle:           yamlConfig.IndexTitle,
		IndexShortcutIconURL: yamlConfig.IndexShortcutIconURL,
	}
	for path, yamlPackage := range yamlConfig.Packages {
		tmplPackage, err := newTmplPackage(yamlConfig.URL, path, yamlPackage)
		if err != nil {
			return nil, err
		}
		tmplConfig.Packages = append(tmplConfig.Packages, tmplPackage)
	}
	for _, yamlIndexLink := range yamlConfig.IndexLinks {
		tmplConfig.IndexLinks = append(
			tmplConfig.IndexLinks,
			&tmplIndexLink{
				URL:  yamlIndexLink.URL,
				Name: yamlIndexLink.Name,
			},
		)
	}
	return tmplConfig, nil
}

func newTmplPackage(configURL string, path string, yamlPackage *yamlPackage) (*tmplPackage, error) {
	tmplPackage := &tmplPackage{
		GoPackage: fmt.Sprintf("%s/%s", configURL, path),
		Path:      path,
		Private:   yamlPackage.Private,
	}
	switch yamlPackage.Type {
	case "github":
		tmplPackage.VCS = "git"
		tmplPackage.IndexURL = fmt.Sprintf("https://github.com/%s/%s", yamlPackage.GithubUser, yamlPackage.GithubRepo)
		tmplPackage.Repository = fmt.Sprintf("https://github.com/%s/%s", yamlPackage.GithubUser, yamlPackage.GithubRepo)
	default:
		return nil, pkgerrors.New("unknown package type", "type", yamlPackage.Type)
	}
	for _, yamlBadge := range yamlPackage.Badges {
		tmplImage := &tmplImage{}
		switch yamlBadge {
		case "circleci":
			switch yamlPackage.Type {
			case "github":
				tmplImage.URL = fmt.Sprintf("https://circleci.com/gh/%s/%s/tree/master", yamlPackage.GithubUser, yamlPackage.GithubRepo)
				tmplImage.Source = fmt.Sprintf("https://circleci.com/gh/%s/%s/tree/master.png", yamlPackage.GithubUser, yamlPackage.GithubRepo)
			default:
				return nil, pkgerrors.New("unknown package type", "type", yamlPackage.Type)
			}
		case "goreportcard":
			tmplImage.URL = fmt.Sprintf("http://goreportcard.com/report/%s", tmplPackage.GoPackage)
			tmplImage.Source = fmt.Sprintf("http://goreportcard.com/badge/%s", tmplPackage.GoPackage)
		case "godoc":
			tmplImage.URL = fmt.Sprintf("https://godoc.org/%s", tmplPackage.GoPackage)
			tmplImage.Source = "http://img.shields.io/badge/GoDoc-Reference-blue.svg"
		case "license-mit":
			switch yamlPackage.Type {
			case "github":
				tmplImage.URL = fmt.Sprintf("https://github.com/%s/%s/blob/master/LICENSE", yamlPackage.GithubUser, yamlPackage.GithubRepo)
			default:
				return nil, pkgerrors.New("unknown package type", "type", yamlPackage.Type)
			}
			tmplImage.Source = "http://img.shields.io/badge/License-MIT-blue.svg"
		default:
			return nil, pkgerrors.New("unknown badge type", "type", yamlBadge)
		}
		tmplPackage.Images = append(tmplPackage.Images, tmplImage)
	}
	return tmplPackage, nil
}
