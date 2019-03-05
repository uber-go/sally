package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	path, clean := TempFile(t, `

url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go

`)
	defer clean()

	config, err := Parse(path)
	assert.NoError(t, err)

	assert.Equal(t, config.Godoc.Host, "godoc.org")
	assert.Equal(t, config.URL, "google.golang.org")

	pkg, ok := config.Packages["grpc"]
	assert.True(t, ok)

	assert.Equal(t, pkg, Package{Repo: "github.com/grpc/grpc-go"})
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
			path, clean := TempFile(t, fmt.Sprintf(`
godoc:
  host: %q
url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
`, tt.give))
			defer clean()

			config, err := Parse(path)
			require.NoError(t, err)

			assert.Equal(t, tt.want, config.Godoc.Host)
			assert.Equal(t, "google.golang.org", config.URL)

			pkg, ok := config.Packages["grpc"]
			assert.True(t, ok)
			assert.Equal(t, Package{Repo: "github.com/grpc/grpc-go"}, pkg)
		})
	}
}

func TestNotAlphabetical(t *testing.T) {
	path, clean := TempFile(t, `

url: google.golang.org
packages:
  grpc:
    repo: github.com/grpc/grpc-go
  atomic:
    repo: github.com/uber-go/atomic

`)
	defer clean()

	_, err := Parse(path)
	if assert.Error(t, err, "YAML configuration is not listed alphabetically") {
		assert.Contains(t, err.Error(), "must be alphabetically ordered")
	}
}
