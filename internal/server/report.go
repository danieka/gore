package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/danieka/gore/internal/reports"
	"github.com/gorilla/mux"
)

var formatContentType map[string]string = map[string]string{
	"csv":  "text/csv",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"json": "application/json",
}

func report(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	parts := strings.Split(name, ".")
	format := "json"
	if len(parts) > 1 {
		format = parts[1]
	}
	report := reports.Reports[parts[0]]
	if report == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Report not found: " + parts[0]))
		return
	}

	output, err := report.Execute(format, r.URL.Query())
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", formatContentType[format])

	_, err = w.Write([]byte(output))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: " + err.Error()))
		return
	}
}
