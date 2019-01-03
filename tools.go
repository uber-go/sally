// +build tools

package main

import (
	_ "github.com/go-bindata/go-bindata/go-bindata"
	_ "github.com/golang/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
