package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"go.uber.org/sally/templates"
)

var (
	indexTemplate = template.Must(
		template.New("index.html").Parse(templates.Index))
	packageTemplate = template.Must(
		template.New("package.html").Parse(templates.Package))
)

// CreateHandler creates a Sally http.Handler
func CreateHandler(config *Config) http.Handler {
	mux := http.NewServeMux()

	pkgs := make([]packageInfo, 0, len(config.Packages))
	for name, pkg := range config.Packages {
		handler := newPackageHandler(config, name, pkg)
		// Double-register so that "/foo"
		// does not redirect to "/foo/" with a 300.
		mux.Handle("/"+name, handler)
		mux.Handle("/"+name+"/", handler)

		pkgs = append(pkgs, packageInfo{
			Desc:       pkg.Desc,
			ImportPath: handler.canonicalURL,
			GitURL:     handler.gitURL,
			GodocHome:  handler.godocHost + "/" + handler.canonicalURL,
		})
	}
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ImportPath < pkgs[j].ImportPath
	})
	mux.Handle("/", &indexHandler{pkgs: pkgs})

	return mux
}

type indexHandler struct {
	pkgs []packageInfo
}

type packageInfo struct {
	Desc       string // package description
	ImportPath string // canonical improt path
	GitURL     string // URL of the Git repository
	GodocHome  string // documentation home URL
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Index handler only supports '/'.
	// ServeMux will call us for any '/foo' that is not a known package.
	if r.Method != http.MethodGet || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct{ Packages []packageInfo }{
		Packages: h.pkgs,
	}
	if err := indexTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type packageHandler struct {
	// Hostname of the godoc server, e.g. "godoc.org".
	godocHost string

	// Name of the package relative to the vanity base URL.
	// For example, "zap" for "go.uber.org/zap".
	name string

	// Path at which the Git repository is hosted.
	// For example, "github.com/uber-go/zap".
	gitURL string

	// Default branch of the Git repository.
	defaultBranch string

	// Canonical import path for the package.
	canonicalURL string
}

func newPackageHandler(cfg *Config, name string, pkg PackageConfig) *packageHandler {
	baseURL := cfg.URL
	if pkg.URL != "" {
		baseURL = pkg.URL
	}
	canonicalURL := fmt.Sprintf("%s/%s", baseURL, name)

	return &packageHandler{
		godocHost:     cfg.Godoc.Host,
		name:          name,
		canonicalURL:  canonicalURL,
		gitURL:        pkg.Repo,
		defaultBranch: pkg.Branch,
	}
}

func (h *packageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Extract the relative path to subpackages, if any.
	//	"/foo/bar" => "/bar"
	//	"/foo" => ""
	relPath := strings.TrimPrefix(r.URL.Path, "/"+h.name)

	data := struct {
		Repo         string
		CanonicalURL string
		GodocURL     string
	}{
		Repo:         h.gitURL,
		CanonicalURL: h.canonicalURL,
		GodocURL:     fmt.Sprintf("https://%s/%s%s", h.godocHost, h.canonicalURL, relPath),
	}
	if err := packageTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
