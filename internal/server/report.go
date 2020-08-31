package server

import (
	"log"
	"net/http"

	"github.com/danieka/gore/internal/reports"
	"github.com/gorilla/mux"
)

func report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	name := vars["name"]
	report := reports.Reports[name]
	output, err := report.Execute("json")
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write([]byte(output))
	if err != nil {
		log.Println(err)
	}
}
