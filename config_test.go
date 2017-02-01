package sally

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, config.URL, "google.golang.org")

	pkg, ok := config.Packages["grpc"]
	assert.True(t, ok)

	assert.Equal(t, pkg, Package{Repo: "github.com/grpc/grpc-go"})
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
