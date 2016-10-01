package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ListenAndServe starts a Sally server
func ListenAndServe(port int, config Config) error {
	router := httprouter.New()
	router.RedirectTrailingSlash = false

	handle, err := createIndexHandle(config)
	if err != nil {
		return err
	}
	router.GET("/", handle)

	for name, p := range config.Packages {
		handle, err := createPackageHandle(pkgViewModel{
			Package: p,
			Name:    name,
			Config:  config,
		})
		if err != nil {
			return err
		}
		router.GET(fmt.Sprintf("/%s", name), handle)
		router.GET(fmt.Sprintf("/%s/*path", name), handle)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		return err
	}

	return nil
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

func createPackageHandle(pvm pkgViewModel) (httprouter.Handle, error) {
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
		t.Execute(w, pvm.NewWithAddlURL(ps.ByName("path")))
	}, nil
}

type pkgViewModel struct {
	Package

	Name    string
	Config  Config
	AddlURI string
}

func (p pkgViewModel) CanonicalURL() string {
	return fmt.Sprintf("%s/%s", p.Config.URL, p.Name)
}

func (p pkgViewModel) GodocURL() string {
	return fmt.Sprintf("https://godoc.org/%s%s", p.CanonicalURL(), p.AddlURI)
}

func (p pkgViewModel) NewWithAddlURL(uri string) pkgViewModel {
	if uri == "" {
		return p
	}
	return pkgViewModel{
		Package: p.Package,
		Name:    p.Name,
		Config:  p.Config,
		AddlURI: uri,
	}
}
