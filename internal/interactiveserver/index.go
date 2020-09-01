package interactiveserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/danieka/gore/internal/reports"
)

var indexTemplate *template.Template

func index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, map[string]interface{}{
		"Reports": reports.Reports,
	})
	if err != nil {
		log.Println(err)
	}
}
