package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	indexTplPath    = "templates/index.tpl"
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

	t, err := template.New(filepath.Base(indexTplPath)).Parse(string(tpl))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, c); err != nil {
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

	t, err := template.New(filepath.Base(packagesTplPath)).Parse(string(tpl))
	if err != nil {
		return err
	}

	for name, pkg := range c.Packages {
		canonicalURL := fmt.Sprintf("%s/%s", c.URL, name)
		tpl := struct {
			Name         string
			CanonicalURL string
			GodocURL     string
			Package
		}{
			Name:         name,
			CanonicalURL: canonicalURL,
			GodocURL:     fmt.Sprintf("https://godoc.org/%s", canonicalURL),
			Package:      pkg,
		}

		buf := new(bytes.Buffer)
		if err := t.Execute(buf, tpl); err != nil {
			return err
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.html", outDir, name), buf.Bytes(), 0644); err != nil {
			return err
		}

		fmt.Println(buf)
	}

	return nil
}
