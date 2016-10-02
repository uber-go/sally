package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	yml := `
url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
`
	path, clean := TempFile(t, yml)
	defer clean()

	config, err := Parse(path)
	assert.NoError(t, err)

	assert.Equal(t, config.URL, "google.golang.org")

	pkg, ok := config.Packages["grpc"]
	assert.True(t, ok)

	assert.Equal(t, pkg, Package{Repo: "github.com/grpc/grpc-go"})
}
