package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// CreateHandler creates a Sally http.Handler
func CreateHandler(config Config) (http.Handler, error) {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	handle, err := createIndexHandle(config)
	if err != nil {
		return router, err
	}
	router.GET("/", handle)

	for name, pkg := range config.Packages {
		handle, err := createPackageHandle(packageViewModel{
			Package: pkg,
			Name:    name,
			Config:  config,
		})
		if err != nil {
			return router, err
		}
		router.GET(fmt.Sprintf("/%s", name), handle)
		router.GET(fmt.Sprintf("/%s/*path", name), handle)
	}

	return router, nil
}

func createIndexHandle(config Config) (httprouter.Handle, error) {
	t, err := template.New("index").Parse(`
<!DOCTYPE html>
<html>
	<body>
		<ul>
			{{ range $key, $value := .Packages }}
				<li>{{ $key }} - {{ $value.Repo }}</li>
			{{ end }}
		</ul>
	</body>
</html>
`)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t.Execute(w, config)
	}, nil
}

func createPackageHandle(pvm packageViewModel) (httprouter.Handle, error) {
	t, err := template.New("package").Parse(`
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="{{ .CanonicalURL }} git https://{{ .Repo }}">
        <meta name="go-source" content="{{ .CanonicalURL }} https://{{ .Repo }} https://{{ .Repo }}/tree/master{/dir} https://{{ .Repo }}/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url={{ .GodocURL }}">
    </head>
    <body>
        Nothing to see here. Please <a href="{{ .GodocURL }}">move along</a>.
    </body>
</html>
`)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t.Execute(w, pvm.NewWithAddlGodocPath(ps.ByName("path")))
	}, nil
}

type packageViewModel struct {
	Package

	Name          string
	Config        Config
	AddlGodocPath string
}

func (p packageViewModel) CanonicalURL() string {
	return fmt.Sprintf("%s/%s", p.Config.URL, p.Name)
}

func (p packageViewModel) GodocURL() string {
	return fmt.Sprintf("https://godoc.org/%s%s", p.CanonicalURL(), p.AddlGodocPath)
}

func (p packageViewModel) NewWithAddlGodocPath(uri string) packageViewModel {
	if uri == "" {
		return p
	}
	return packageViewModel{
		Package:       p.Package,
		Name:          p.Name,
		Config:        p.Config,
		AddlGodocPath: uri,
	}
}
