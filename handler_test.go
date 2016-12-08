package sally

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPage(t *testing.T) {
	pathToPage := map[string][]byte{
		"foo":        []byte("1"),
		"bar":        []byte("2"),
		"foo.master": []byte("3"),
	}
	checkPage(t, pathToPage, "1", "foo")
	checkPage(t, pathToPage, "2", "bar")
	checkPage(t, pathToPage, "3", "foo.master")
	checkPage(t, pathToPage, "1", "foo/bar")
	checkPage(t, pathToPage, "2", "bar/bar")
	checkPage(t, pathToPage, "3", "foo.master/foo")
}

func checkPage(t *testing.T, pathToPage map[string][]byte, pageString string, path string) {
	require.Equal(t, []byte(pageString), getPage(pathToPage, path))
}
