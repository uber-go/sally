// +build tools

package main

import (
	_ "github.com/golang/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
