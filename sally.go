/*
Package sally defines functionality for an http.Handler that handles go import path redirection.
*/
package sally // import "go.uber.org/sally"

import "net/http"

// NewHandler returns a new http.Handler.
func NewHandler(goTmplFilePath string, indexTmplFilePath string, configFilePath string) (http.Handler, error) {
	pathToPage, err := generatePathToPage(goTmplFilePath, indexTmplFilePath, configFilePath)
	if err != nil {
		return nil, err
	}
	return newHandler(pathToPage), nil
}
