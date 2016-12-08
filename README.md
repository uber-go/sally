# sally

A tiny HTTP server for supporting custom Golang import paths

## Installation

`make install`

## Usage

Create a YAML file with the following structure:

```yaml
url: go.uber.org
index_title: uber-go
packages:
  thriftrw:
    type: github
    github_user: thriftrw
    github_repo: thriftrw-go
    badges:
      - godoc
  yarpc:
    type: github
    github_user: yarpc
    github_repo: yarpc-go
    badges:
      - godoc
```

Then run Sally to start the HTTP server:

```
$ make launch
```
