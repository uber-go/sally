# sally

A tiny HTTP server for supporting custom Golang import paths

## Installation

`go get go.uber.org/sally`

## Usage

Create a YAML file with the following structure:

```yaml
# This optional section configures godoc documentation linking.
godoc:
  # Instance of godoc server used for documentation links. Defaults to pkg.go.dev.
  host: pkg.go.dev

url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
```

Then run Sally to start the HTTP server:

```
$ sally -yml site.yaml -port 5000
```
