package util

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/minecrafter/sage/bindata"
)

var (
	fourOhFour = bindata.MustAsset("templates/404.html")
	base       *template.Template // populated later
)

func init() {
	files := []string{"templates/footer.html", "templates/header.html", "templates/internal_error.html",
		"templates/main.html"}
	var buf bytes.Buffer
	for _, file := range files {
		buf.Write(bindata.MustAsset(file))
	}
	base = template.Must(template.New("base").Parse(buf.String()))
}

func mustTemplateAsset(name string) *template.Template {
	return template.Must(template.New(name).Parse(string(bindata.MustAsset(
		"templates/" + name + ".html"))))
}

func Do404(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	w.Write(fourOhFour)
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
		Error: err.Error(),
	})
}

func DoMain(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	if err := base.ExecuteTemplate(w, "main", nil); err != nil {
		log.Println("Error while rendering", err)
	}
}
