package interactiveserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/danieka/gore/internal/reports"
)

var indexTemplate *template.Template

func init() {
	var err error
	indexTemplate, err = template.New("").ParseFiles("internal/interactiveserver/templates/index.html", "internal/interactiveserver/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.ExecuteTemplate(w, "base", map[string]interface{}{
		"Reports": reports.Reports,
	})
	if err != nil {
		log.Println(err)
	}
}
