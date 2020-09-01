package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Start the web server
func Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/reports/", index)
	r.HandleFunc("/reports/{name}", report)

	http.Handle("/", r)

	fmt.Println("Starting server on :16772")
	err := http.ListenAndServe(":16772", nil)
	return err
}
