package interactiveserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/danieka/gore/internal/reports"
	"github.com/gorilla/mux"
)

var reportTemplate *template.Template

func report(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	report := reports.Reports[name]
	output, err := report.Execute("json")
	if err != nil {
		log.Println(err)
	}
	err = reportTemplate.Execute(w, map[string]interface{}{
		"Report": report,
		"Output": output,
	})
	if err != nil {
		log.Println(err)
	}
}
