//go:build tools
// +build tools

package main

import (
	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
