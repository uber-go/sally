package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func GenerateSite(config *Config, output string) (err error) {
	log.Printf("Generating static site to %s ...", output)
	log.Printf("Generating index.html ...")
	err = os.MkdirAll(output, 0755)
	if err != nil {
		return
	}
	var f *os.File
	f, err = os.OpenFile(path.Join(output, "index.html"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	if err = indexTemplate.Execute(f, config); err != nil {
		return
	}
	_ = f.Close()

	for name, pkg := range config.Packages {
		log.Printf("Generating %s ...", name)
		err = os.MkdirAll(path.Dir(path.Join(output, name)), 0755)
		if err != nil {
			return
		}
		canonicalURL := fmt.Sprintf("%s/%s", config.URL, name)
		data := struct {
			Repo         string
			Branch       string
			CanonicalURL string
			GodocURL     string
		}{
			Repo:         pkg.Repo,
			Branch:       pkg.Branch,
			CanonicalURL: canonicalURL,
			GodocURL:     fmt.Sprintf("https://%s/%s", config.Godoc.Host, canonicalURL),
		}
		f, err = os.OpenFile(path.Join(output, name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		if err = packageTemplate.Execute(f, data); err != nil {
			return
		}
		_ = f.Close()
	}
	return
}
