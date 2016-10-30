// Copyright 2016 Package Savant team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"net/http"
	"strings"

	"github.com/minecrafter/package-savant/repository/maven"
	"github.com/minecrafter/package-savant/util"
	"github.com/pkg/errors"
)

type RepoHTTPHandler struct {
	Repositories map[string]maven.MavenRetrieveHandler
}

func (h RepoHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// embellish ourselves
	w.Header().Add("Server", "Package Savant")

	path := r.URL.EscapedPath()
	if path == "/ping" {
		w.Write([]byte("ok"))
		return
	}

	if path == "/" {
		// get repository names and total package count
		packageCount := 0
		repositoryNames := make([]string, len(h.Repositories))
		i := 0
		for key, server := range h.Repositories {
			packages, err := server.MetadataStore.GetAllIDs()
			if err != nil {
				// Skip this
				continue
			}
			repositoryNames[i] = key
			i++
			packageCount += len(*packages)
		}
		util.DoMain(w, repositoryNames, packageCount)
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
