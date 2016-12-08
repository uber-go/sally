package main

import (
	"go.pedge.io/lion/env"
	"go.pedge.io/pkg/http"
	"go.uber.org/sally"
)

type appEnv struct {
	GoTemplateFilePath    string `env:"GO_TEMPLATE,required"`
	IndexTemplateFilePath string `env:"INDEX_TEMPLATE,required"`
	ConfigFilePath        string `env:"CONFIG,required"`
}

func main() {
	envlion.Main(func(appEnvObj interface{}) error { return do(appEnvObj.(*appEnv)) }, &appEnv{})
}

func do(appEnv *appEnv) error {
	handler, err := sally.NewHandler(
		appEnv.GoTemplateFilePath,
		appEnv.IndexTemplateFilePath,
		appEnv.ConfigFilePath,
	)
	if err != nil {
		return err
	}
	return pkghttp.GetAndListenAndServe(handler)
}
