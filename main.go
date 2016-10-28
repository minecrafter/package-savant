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

package main

import (
	"log"
	"net/http"

	"github.com/minecrafter/sage/repository/maven"
	"github.com/minecrafter/sage/repository/store"
	"github.com/minecrafter/sage/server"
)

func main() {
	conf, err := server.LoadConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	configuredRepositories := make(map[string]maven.MavenRetrieveHandler)
	for key, value := range conf.Repositories {
		configuredRepositories[key] = maven.NewMavenRetrieveHandler("/"+key,
			store.NewLocalMetadataStore(value.Metadata.Path),
			store.NewLocalPackageStore(value.Storage.Path))
	}

	handler := server.RepoHTTPHandler{Repositories: configuredRepositories}
	if conf.SSL.Enabled {
		err = http.ListenAndServeTLS(conf.Listen, conf.SSL.CertFile, conf.SSL.KeyFile, handler)
	} else {
		err = http.ListenAndServe(conf.Listen, handler)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
