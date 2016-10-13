package server

import (
	"net/http"
	"strings"

	"github.com/minecrafter/sage/repository/maven"
)

type RepoHTTPHandler struct {
	Repositories map[string]maven.MavenRetrieveHandler
}

func (h RepoHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// embellish ourselves
	w.Header().Add("Server", "Sage")
	w.Header().Add("X-Sage-Version", "0.1")
	w.Header().Add("X-Try-Sage", "http://sage-repo.org")

	path := r.URL.EscapedPath()
	if path == "/" {
		w.Write([]byte("ok"))
		return
	}

	firstSlash := strings.Index(path[1:], "/")
	if firstSlash < 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("repository not found"))
		return
	}
	repoName := path[1 : firstSlash+1]
	repoServer, exists := h.Repositories[repoName]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("repository not found"))
		return
	}

	if r.Method == "GET" {
		repoServer.GetMavenFile(w, r)
	} else if r.Method == "PUT" {
		repoServer.PutMavenFile(w, r, path[firstSlash+1:])
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
