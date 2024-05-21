package main

import (
	"cmp"
	"embed"
	"errors"
	"html/template"
	"net/http"
	"path"
	"slices"
	"strings"
)

var (
	//go:embed templates/*.html
	templateFiles embed.FS

	_templates = template.Must(template.ParseFS(templateFiles, "templates/*.html"))
)

// CreateHandler builds a new handler with the provided package configuration,
// and templates. The templates object must contain the following: index.html,
// package.html, and 404.html. The returned handler provides the following
// endpoints:
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
func CreateHandler(config *Config, templates *template.Template) (http.Handler, error) {
	indexTemplate := templates.Lookup("index.html")
	if indexTemplate == nil {
		return nil, errors.New("template index.html is missing")
	}

	notFoundTemplate := templates.Lookup("404.html")
	if notFoundTemplate == nil {
		return nil, errors.New("template 404.html is missing")
	}

	packageTemplate := templates.Lookup("package.html")
	if packageTemplate == nil {
		return nil, errors.New("template package.html is missing")
	}

	mux := http.NewServeMux()
	pkgs := make([]*sallyPackage, 0, len(config.Packages))
	for name, pkg := range config.Packages {
		baseURL := config.URL
		if pkg.URL != "" {
			// Package-specific override for the base URL.
			baseURL = pkg.URL
		}
		modulePath := path.Join(baseURL, name)

		docURL := pkg.DocURL
		if docURL == "" {
			docURL = "https://" + path.Join(config.Godoc.Host, modulePath)
		}

		docBadge := pkg.DocBadge
		if docBadge == "" {
			docBadge = "//pkg.go.dev/badge/" + modulePath + ".svg"
		}

		pkg := &sallyPackage{
			Name:       name,
			Desc:       pkg.Desc,
			ModulePath: modulePath,
			DocURL:     docURL,
			DocBadge:   docBadge,
			VCS:        pkg.VCS,
			RepoURL:    pkg.Repo,
		}
		pkgs = append(pkgs, pkg)

		// Double-register so that "/foo"
		// does not redirect to "/foo/" with a 300.
		handler := &packageHandler{pkg: pkg, template: packageTemplate}
		mux.Handle("/"+name, handler)
		mux.Handle("/"+name+"/", handler)
	}

	mux.Handle("/", newIndexHandler(pkgs, indexTemplate, notFoundTemplate))
	return requireMethod(http.MethodGet, mux), nil
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

	// URL at which documentation badge image can be found.
	DocBadge string

	// Version control system used by the package.
	VCS string

	// URL at which the repository is hosted.
	RepoURL string
}

type indexHandler struct {
	pkgs             []*sallyPackage // sorted by name
	indexTemplate    *template.Template
	notFoundTemplate *template.Template
}

var _ http.Handler = (*indexHandler)(nil)

func newIndexHandler(pkgs []*sallyPackage, indexTemplate, notFoundTemplate *template.Template) *indexHandler {
	slices.SortFunc(pkgs, func(a, b *sallyPackage) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return &indexHandler{
		pkgs:             pkgs,
		indexTemplate:    indexTemplate,
		notFoundTemplate: notFoundTemplate,
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
		serveHTML(w, http.StatusNotFound, h.notFoundTemplate, struct{ Path string }{
			Path: path,
		})
		return
	}

	serveHTML(w, http.StatusOK, h.indexTemplate, struct{ Packages []*sallyPackage }{
		Packages: h.pkgs[start:end],
	})
}

type packageHandler struct {
	pkg      *sallyPackage
	template *template.Template
}

var _ http.Handler = (*packageHandler)(nil)

func (h *packageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the relative path to subpackages, if any.
	//      "/foo/bar" => "/bar"
	//      "/foo" => ""
	relPath := strings.TrimPrefix(r.URL.Path, "/"+h.pkg.Name)

	serveHTML(w, http.StatusOK, h.template, struct {
		ModulePath string
		VCS        string
		RepoURL    string
		DocURL     string
	}{
		ModulePath: h.pkg.ModulePath,
		VCS:        h.pkg.VCS,
		RepoURL:    h.pkg.RepoURL,
		DocURL:     h.pkg.DocURL + relPath,
	})
}

func descends(from, to string) bool {
	return to == from || (strings.HasPrefix(to, from) && to[len(from)] == '/')
}

func serveHTML(w http.ResponseWriter, status int, template *template.Template, data interface{}) {
	if status >= 400 {
		w.Header().Set("Cache-Control", "no-cache")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	err := template.Execute(w, data)
	if err != nil {
		// The status has already been sent, so we cannot use [http.Error] - otherwise
		// we'll get a superfluous call warning. The other option is to execute the template
		// to a temporary buffer, but memory.
		_, _ = w.Write([]byte(err.Error()))
	}
}
