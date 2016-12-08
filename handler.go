package sally

import (
	"net/http"
	"strings"
)

type handler struct {
	pathToPage map[string][]byte
}

func newHandler(pathToPage map[string][]byte) *handler {
	return &handler{pathToPage}
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(strings.TrimSuffix(request.URL.Path, "/"), "/")
	if page := getPage(h.pathToPage, path); page != nil {
		_, _ = responseWriter.Write(page)
		return
	}
	http.NotFound(responseWriter, request)
}

func getPage(pathToPage map[string][]byte, path string) []byte {
	page, ok := pathToPage[path]
	if ok {
		return page
	}
	split := strings.Split(path, "/")
	lenSplit := len(split)
	if lenSplit > 1 {
		for i := lenSplit - 1; i > 0; i-- {
			path = strings.Join(split[0:i], "/")
			page, ok = pathToPage[path]
			if ok {
				return page
			}
		}
	}
	return nil
}
