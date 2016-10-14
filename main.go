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
