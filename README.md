# sally

A canonical import path static site generator for Go

## Installation

`go get go.uber.org/sally`

## Usage

Create a YAML file with the following structure:

```yaml
url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
```

Then run Sally to start the HTTP server:

```
$ sally -yml site.yaml -port 5000
```
