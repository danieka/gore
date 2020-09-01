package reports

var defaultReports = map[string]string{
	"csv": `{{$cols := .Cols }}{{ range $index, $col := .Cols}}{{if $index}},{{end}}{{$col}}{{end}}{{ range $row := .Rows }}
{{ range $index, $col := $cols}}{{if $index}},{{end}}"{{index $row $col}}"{{end}}{{end}}`,
}
