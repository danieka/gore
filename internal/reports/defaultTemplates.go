package reports

var defaultReports = map[string]string{
	"csv": `{{$cols := .Cols }}{{ range $index, $col := .Cols}}{{if $index}},{{end}}"{{$col}}"{{end}}{{ range $row := .Rows }}
{{ range $index, $col := $cols}}{{if $index}},{{end}}"{{index $row $col}}"{{end}}{{end}}`,
	"json": `[
		{{$cols := .Cols }}
		{{ range $rowIdx, $row := .Rows }}
		{{if $rowIdx}},{{end}}
		{
			{{ range $colIdx, $col := $cols}}
				{{if $colIdx}},{{end}}
				"{{ $col }}": "{{ index $row $col }}"
			{{ end }}
		}
		{{ end }}
		]`,
}
