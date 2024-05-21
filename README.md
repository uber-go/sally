# sally

sally is a small HTTP service you can host
to serve vanity import paths for Go modules.

## Installation

To build sally from source, use:

```bash
go install go.uber.org/sally@latest
```

Alternatively, get a pre-built Docker image from
https://github.com/uber-go/sally/pkgs/container/sally.

## Usage

Create a YAML file with the following structure:

```yaml
# sally.yaml

# Configures documentation linking.
# Optional.
godoc:
  # Host for the Go documentation server.
  # Defaults to pkg.go.dev.
  host: pkg.go.dev

# Base URL for your package site.
# If you want your modules available under "example.com",
# specify example.com here.
# This field is required.
url: go.uber.org

# Collection of packages under example.com
# and their Git repositories.
packages:

  # The key is the name of the package following the base URL.
  # For example, if you want to make a package available at
  # "example.com/foo", you'd specify "foo" here.
  zap:
    # Path to the Git repository.
    #
    # This field is required.
    repo: github.com/uber-go/zap

    # Optional description of the package.
    description: A fast, structured-logging library.

    # Alternative base URL instead of the value configured at the top-level.
    # This is useful if the same sally instance is
    # hosted behind multiple base URLs.
    #
    # Defaults to the value of the top-level url field.
    url: example.com

    # Optional URL to the package's documentation.
    #
    # Defaults to the documentation site at pkg.go.dev with the package's
    # module path appended.
    doc_url: example.com/go-pkg/docs/zap

    # Optional URL to the badge image which appears in the index page.
    #
    # Defaults to the badge image at pkg.go.dev, using the package's module
    # path followed by .svg as the filename.
    doc_badge: example.com/go-pkg/badge/zap
```

Run sally like so:

```shell
$ sally
```

This will read from sally.yaml and serve on port 8080.
To use a different port and/or configuration file,
use the `-yml` and `-port` flags.

```
$ sally -yml site.yaml -port 5000
```

### Custom Templates

You can provide your own custom templates. For this, create a directory with `.html`
templates and provide it via the `-templates` flag. You only need to provide the
templates you want to override. See [templates](./templates/) for the available
templates.
