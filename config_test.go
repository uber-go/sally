package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	path := TempFile(t, `

url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
    branch: main
    vcs: svn

`)

	config, err := Parse(path)
	assert.NoError(t, err)

	assert.Equal(t, config.Godoc.Host, "pkg.go.dev")
	assert.Equal(t, config.URL, "google.golang.org")

	pkg, ok := config.Packages["grpc"]
	assert.True(t, ok)
	assert.Equal(t, PackageConfig{Repo: "github.com/grpc/grpc-go", VCS: "svn"}, pkg)
}

func TestParsePackageLevelURL(t *testing.T) {
	path := TempFile(t, `

url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
    url: go.uber.org

`)

	config, err := Parse(path)
	assert.NoError(t, err)

	pkg, ok := config.Packages["grpc"]
	assert.True(t, ok)
	assert.Equal(t, pkg.URL, "go.uber.org")
}

func TestParseGodocServer(t *testing.T) {
	tests := []struct {
		give string
		want string
	}{
		{"example.com", "example.com"},
		{"example.com/", "example.com"},
		{"http://example.com/", "example.com"},
		{"https://example.com/", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			path := TempFile(t, fmt.Sprintf(`
godoc:
  host: %q
url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
`, tt.give))

			config, err := Parse(path)
			require.NoError(t, err)

			assert.Equal(t, tt.want, config.Godoc.Host)
			assert.Equal(t, "google.golang.org", config.URL)

			pkg, ok := config.Packages["grpc"]
			assert.True(t, ok)
			assert.Equal(t, PackageConfig{Repo: "github.com/grpc/grpc-go", VCS: "git"}, pkg)
		})
	}
}
