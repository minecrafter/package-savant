package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/minecrafter/sage/repository/maven"
	"github.com/minecrafter/sage/util"
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
	if path == "/ping" {
		w.Write([]byte("ok"))
		return
	}

	if path == "/" {
		util.DoMain(w)
		return
	}

	if path == "/fuck" {
		util.DoSpecificError(w, errors.New("You dun goofed"))
		return
	}

	firstSlash := strings.Index(path[1:], "/")
	if firstSlash == -1 {
		// Not handled explicitly earlier, so bomb out.
		util.Do404(w)
		return
	}
	repoName := path[1 : firstSlash+1]
	repoServer, exists := h.Repositories[repoName]
	if !exists {
		util.Do404(w)
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
