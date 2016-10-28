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

package util

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"fmt"

	"github.com/minecrafter/sage/bindata"
)

var (
	base *template.Template // populated later
)

func init() {
	files := []string{"templates/footer.html", "templates/header.html", "templates/internal_error.html",
		"templates/main.html", "templates/404.html"}
	var concatenatedTemplate bytes.Buffer
	for _, file := range files {
		concatenatedTemplate.Write(bindata.MustAsset(file))
	}
	base = template.Must(template.New("base").Parse(concatenatedTemplate.String()))
}

func mustTemplateAsset(name string) *template.Template {
	return template.Must(template.New(name).Parse(string(bindata.MustAsset(
		"templates/" + name + ".html"))))
}

func Do404(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	base.ExecuteTemplate(w, "four_oh_four", nil)
}

func DoGenericError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
	base.ExecuteTemplate(w, "internal_error", struct {
		Error string
	}{})
}

func DoSpecificError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
	base.ExecuteTemplate(w, "internal_error", struct {
		Error string
	}{
		Error: fmt.Sprintf("%+v", err),
	})
}

func DoMain(w http.ResponseWriter, repositories []string, totalCount int) {
	w.Header().Set("Content-Type", "text/html")
	data := struct {
		RepositoryCount int
		Repositories    []string
		PackageCount    int
	}{
		RepositoryCount: len(repositories),
		Repositories:    repositories,
		PackageCount:    totalCount,
	}
	if err := base.ExecuteTemplate(w, "main", data); err != nil {
		log.Println("Error while rendering", err)
	}
}
