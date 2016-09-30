package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchURL_Private(t *testing.T) {
	p := makePkgTpl("github.com/uber-go/fake", true)

	imp := p.FetchURL()
	assert.Equal(t, "git@github.com:uber-go/fake", imp, "Expected correct git SSH url")
}

func TestFetchURL_Public(t *testing.T) {
	p := makePkgTpl("github.com/uber-go/sally", false)

	imp := p.FetchURL()
	assert.Equal(t, "https://github.com/uber-go/sally", imp, "Expected correct https")
}

func makePkgTpl(repo string, private bool) packageTpl {
	return packageTpl{
		Package: Package{
			Private: private,
			Repo:    repo,
		},
	}
}
