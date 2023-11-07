package main

import (
	"cmp"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"slices"
	"strings"
)

var (
	//go:embed templates/*.html
	templateFiles embed.FS

	templates = template.Must(template.ParseFS(templateFiles, "templates/*.html"))
)

// CreateHandler builds a new handler
// with the provided package configuration.
// The returned handler provides the following endpoints:
//
//	GET /
//		Index page listing all packages.
//	GET /<name>
//	       Package page for the given package.
//	GET /<dir>
//		Page listing packages under the given directory,
//		assuming that there's no package with the given name.
//	GET /<name>/<subpkg>
//		Package page for the given subpackage.
func CreateHandler(config *Config) http.Handler {
	mux := http.NewServeMux()
	pkgs := make([]*sallyPackage, 0, len(config.Packages))
	for name, pkg := range config.Packages {
		baseURL := config.URL
		if pkg.URL != "" {
			// Package-specific override for the base URL.
			baseURL = pkg.URL
		}
		modulePath := path.Join(baseURL, name)
		docURL := "https://" + path.Join(config.Godoc.Host, modulePath)

		pkg := &sallyPackage{
			Name:       name,
			Desc:       pkg.Desc,
			ModulePath: modulePath,
			DocURL:     docURL,
			VCS:        pkg.VCS,
			RepoURL:    pkg.Repo,
		}
		pkgs = append(pkgs, pkg)

		// Double-register so that "/foo"
		// does not redirect to "/foo/" with a 300.
		handler := &packageHandler{Pkg: pkg}
		mux.Handle("/"+name, handler)
		mux.Handle("/"+name+"/", handler)
	}

	mux.Handle("/", newIndexHandler(pkgs))
	return requireMethod(http.MethodGet, mux)
}

func requireMethod(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

type sallyPackage struct {
	// Name of the package.
	//
	// This is the part after the base URL.
	Name string

	// Canonical import path for the package.
	ModulePath string

	// Description of the package, if any.
	Desc string

	// URL at which documentation for the package can be found.
	DocURL string

	// Version control system used by the package.
	VCS string

	// URL at which the repository is hosted.
	RepoURL string
}

type indexHandler struct {
	pkgs []*sallyPackage // sorted by name
}

var _ http.Handler = (*indexHandler)(nil)

func newIndexHandler(pkgs []*sallyPackage) *indexHandler {
	slices.SortFunc(pkgs, func(a, b *sallyPackage) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return &indexHandler{
		pkgs: pkgs,
	}
}

func (h *indexHandler) rangeOf(path string) (start, end int) {
	if len(path) == 0 {
		return 0, len(h.pkgs)
	}

	// If the packages are sorted by name,
	// we can scan adjacent packages to find the range of packages
	// whose name descends from path.
	start, _ = slices.BinarySearchFunc(h.pkgs, path, func(pkg *sallyPackage, path string) int {
		return cmp.Compare(pkg.Name, path)
	})

	for idx := start; idx < len(h.pkgs); idx++ {
		if !descends(path, h.pkgs[idx].Name) {
			// End of matching sequences.
			// The next path is not a descendant of path.
			return start, idx
		}
	}

	// All packages following start are descendants of path.
	// Return the rest of the packages.
	return start, len(h.pkgs)
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/")
	start, end := h.rangeOf(path)

	// If start == end, then there are no packages
	if start == end {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no packages found under path: %v\n", path)
		return
	}

	err := templates.ExecuteTemplate(w, "index.html",
		struct{ Packages []*sallyPackage }{
			Packages: h.pkgs[start:end],
		})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type packageHandler struct {
	Pkg *sallyPackage
}

var _ http.Handler = (*packageHandler)(nil)

func (h *packageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the relative path to subpackages, if any.
	//      "/foo/bar" => "/bar"
	//      "/foo" => ""
	relPath := strings.TrimPrefix(r.URL.Path, "/"+h.Pkg.Name)

	err := templates.ExecuteTemplate(w, "package.html", struct {
		ModulePath string
		VCS        string
		RepoURL    string
		DocURL     string
	}{
		ModulePath: h.Pkg.ModulePath,
		VCS:        h.Pkg.VCS,
		RepoURL:    h.Pkg.RepoURL,
		DocURL:     h.Pkg.DocURL + relPath,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func descends(from, to string) bool {
	return to == from || (strings.HasPrefix(to, from) && to[len(from)] == '/')
}
