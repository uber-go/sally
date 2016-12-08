package sally

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

var (
	indexPaths = []string{
		"",
		"index.html",
		"index.htm",
		"default.htm",
	}
)

func generatePathToPage(goTmplFilePath string, indexTmplFilePath string, configFilePath string) (map[string][]byte, error) {
	goTmpl, err := template.ParseFiles(goTmplFilePath)
	if err != nil {
		return nil, err
	}
	indexTmpl, err := template.ParseFiles(indexTmplFilePath)
	if err != nil {
		return nil, err
	}
	yamlData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	yamlConfig, err := newYAMLConfig(yamlData)
	if err != nil {
		return nil, err
	}
	tmplConfig, err := newTmplConfig(yamlConfig)
	if err != nil {
		return nil, err
	}
	return generatePathToPageParsed(goTmpl, indexTmpl, tmplConfig)
}

func generatePathToPageParsed(goTmpl *template.Template, indexTmpl *template.Template, tmplConfig *tmplConfig) (map[string][]byte, error) {
	pathToPage := make(map[string][]byte)
	for _, tmplPackage := range tmplConfig.Packages {
		page, err := generatePage(goTmpl, tmplPackage)
		if err != nil {
			return nil, err
		}
		pathToPage[tmplPackage.Path] = page
	}
	indexPage, err := generatePage(indexTmpl, tmplConfig)
	if err != nil {
		return nil, err
	}
	for _, indexPath := range indexPaths {
		pathToPage[indexPath] = indexPage
	}
	return pathToPage, nil
}

func generatePage(t *template.Template, data interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := t.Execute(buffer, data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
