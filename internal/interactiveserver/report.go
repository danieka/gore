package interactiveserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/danieka/gore/internal/reports"
	"github.com/gorilla/mux"
)

var reportTemplate *template.Template

func init() {
	var err error
	reportTemplate, err = template.New("").ParseFiles("internal/interactiveserver/templates/report.html", "internal/interactiveserver/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
}

func report(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	report := reports.Reports[name]
	output, err := report.Execute()
	if err != nil {
		log.Println(err)
	}
	err = reportTemplate.ExecuteTemplate(w, "base", map[string]interface{}{
		"Report": report,
		"Output": output,
	})
	if err != nil {
		log.Println(err)
	}
}
