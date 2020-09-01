package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/danieka/gore/internal/reports"
)

var templateStr = `
{{$first := true}}
[
	{{ $host := .Host}}
{{range $key, $value := .Reports}}
    {{if $first}}
        {{$first = false}}
    {{else}}
        ,
	{{end}}
	{
		"name": "{{$value.Info.Name}}",
		"url": "{{$host}}/reports/{{$value.Info.Name}}"
	}
{{end}}
]
`

var indexTemplate *template.Template

func init() {
	var err error
	indexTemplate, err = template.New("").Parse(templateStr)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := indexTemplate.ExecuteTemplate(w, "", map[string]interface{}{
		"Reports": reports.Reports,
		"Host":    "http://localhost:16772",
	})
	if err != nil {
		log.Println(err)
	}
}
