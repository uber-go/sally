package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

const (
	indexTpl     = "index.tpl"
	indexTplPath = "templates/index.tpl"

	packagesTpl     = "package.tpl"
	packagesTplPath = "templates/package.tpl"
)

// Write takes a Config and produces a static html site to outDir
func Write(c Config, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	if err := writeIndex(c, outDir); err != nil {
		return err
	}
	if err := writePackages(c, outDir); err != nil {
		return err
	}
	return nil
}

func writeIndex(c Config, outDir string) error {
	tpl, err := Asset(indexTplPath)
	if err != nil {
		return err
	}

	t, err := template.New(indexTpl).Parse(string(tpl))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/index.html", outDir), buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}

func writePackages(c Config, outDir string) error {
	tpl, err := Asset(packagesTplPath)
	if err != nil {
		return err
	}

	t, err := template.New(packagesTpl).Parse(string(tpl))
	if err != nil {
		return err
	}

	for name, pkg := range c.Packages {
		tpl := struct {
			CanonicalURL string
			Name         string
			Package
		}{
			CanonicalURL: fmt.Sprintf("%s/%s", c.URL, name),
			Name:         name,
			Package:      pkg,
		}

		buf := new(bytes.Buffer)
		err = t.Execute(buf, tpl)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.html", outDir, name), buf.Bytes(), 0644)
		if err != nil {
			return err
		}

		fmt.Println(buf)
	}

	return nil
}
