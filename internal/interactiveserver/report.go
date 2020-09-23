package interactiveserver

import (
	"fmt"
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

	format := "json"
	query := r.URL.Query()
	val, present := query["format"]
	if present {
		format = val[0]
	}
	output, err := report.Execute(format)
	var errorString string
	if err != nil {
		errorString = fmt.Sprintf("Error executing forma %s: %s", format, err)
	}
	err = reportTemplate.ExecuteTemplate(w, "base", map[string]interface{}{
		"Report": report,
		"Output": output,
		"Format": format,
		"Error":  errorString,
	})
	if err != nil {
		log.Println(err)
	}
}
