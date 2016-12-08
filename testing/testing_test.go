package testing

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/sally"
)

func TestBasic(t *testing.T) {
	checkStandardPath(t, "thriftrw", "git", "https://github.com/thriftrw/thriftrw-go", "go.uber.org/thriftrw")
	checkStandardPath(t, "thriftrw/foo", "git", "https://github.com/thriftrw/thriftrw-go", "go.uber.org/thriftrw")
	checkStandardPath(t, "yarpc", "git", "https://github.com/yarpc/yarpc-go", "go.uber.org/yarpc")
}

func TestIndex(t *testing.T) {
	urls := map[string]bool{
		"go.uber.org/thriftrw": true,
		"go.uber.org/yarpc":    true,
	}
	statusCode, body := getStatusCodeAndBody(t, "")
	require.Equal(t, http.StatusOK, statusCode)
	split := strings.Split(body, "\n")
	getUrls := make(map[string]bool)
	for _, url := range split {
		if url != "" {
			getUrls[url] = true
		}
	}
	require.Equal(t, urls, getUrls)
}

func checkStandardPath(t *testing.T, path string, vcs string, repository string, url string) {
	checkResponse(t, path, http.StatusOK, fmt.Sprintf("%s %s %s\n", vcs, repository, url))
}

func checkResponse(t *testing.T, path string, statusCode int, body string) {
	getStatusCode, getBody := getStatusCodeAndBody(t, path)
	require.Equal(t, statusCode, getStatusCode)
	if statusCode == http.StatusOK {
		require.Equal(t, body, getBody)
	}
}

func getStatusCodeAndBody(t *testing.T, path string) (int, string) {
	handler, err := newTestHandler()
	require.NoError(t, err)
	server := httptest.NewServer(handler)
	defer server.Close()
	response, err := http.Get(fmt.Sprintf("%s/%s", server.URL, path))
	require.NoError(t, err)
	body := ""
	if response.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)
		body = string(data)
		require.NoError(t, response.Body.Close())
	}
	return response.StatusCode, body
}

func newTestHandler() (http.Handler, error) {
	return sally.NewHandler("go.html", "index.html", "config.yaml")
}
