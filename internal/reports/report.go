package reports

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/danieka/gore/internal/sources"

	"gopkg.in/yaml.v2"
)

// Parameter for the sql query
type Parameter struct {
	Name       string
	ColumnName string `yaml:"columnName"`
}

// ReportInfo contains all metadata about a report
type ReportInfo struct {
	Name       string
	Parameters []Parameter
}

// ReportSource is the data source for the report
type ReportSource struct {
	sourceName string
	query      string
}

// ReportOutput contains info on the output
type ReportOutput struct {
	Format   string
	template string
}

// Report is the main struct for a report
type Report struct {
	Info    ReportInfo
	source  ReportSource
	Outputs map[string]ReportOutput
}

var columnNames []string = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U"}

// Execute the report and return the output
func (r *Report) Execute(format string, arguments url.Values) (s string, err error) {
	source := sources.Sources[r.source.sourceName]
	if source == nil {
		return "", fmt.Errorf("Unable to find source %s", r.source.sourceName)
	}

	query := r.source.query
	var wheres []string
	var params []interface{}
	for _, param := range r.Info.Parameters {
		inputArgument := arguments.Get(param.Name)
		if inputArgument != "" {
			wheres = append(wheres, fmt.Sprintf("%s in (?)", param.ColumnName))
			params = append(params, inputArgument)
		}
	}

	if len(wheres) > 0 {
		query += "\nWHERE "
		for _, clause := range wheres {
			query += clause + " AND "
		}
		query += "TRUE"
	}

	cols, rows, err := sources.Sources[r.source.sourceName].Execute(query, params...)
	if err != nil {
		return "", fmt.Errorf("Error executing query %s\nComplete query was %s", err, query)
	}

	output, ok := r.Outputs[format]
	if !ok {
		return "", fmt.Errorf("Unable to find output format %s in outputs %v", format, r.Outputs)
	}

	var buf bytes.Buffer

	if format == "xlsx" {
		f := excelize.NewFile()
		for i, v := range cols {
			f.SetCellValue("Sheet1", columnNames[i]+"1", v)
		}
		for rowIndex, row := range rows {
			for columnIndex, colName := range cols {
				f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", columnNames[columnIndex], rowIndex+2), row[colName])
			}
		}
		// Save xlsx file by the given path.
		bufP, err := f.WriteToBuffer()
		if err != nil {
			log.Fatal(err)
		}
		buf = *bufP
	} else {
		var tmpl *template.Template
		tmpl, err = template.New("").Parse(output.template)
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.ExecuteTemplate(&buf, "", map[string]interface{}{
			"Rows": rows,
			"Cols": cols,
		})
		if err != nil {
			log.Println(err)
		}
	}

	return buf.String(), err
}

// Reports contain all reports that have been read
var Reports map[string]*Report = make(map[string]*Report)

func parseInfo(scanner *bufio.Scanner) (info ReportInfo, err error) {
	var data string
L:
	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "</info>":
			break L
		default:
			data = data + text + "\n"
		}
	}
	err = yaml.Unmarshal([]byte(data), &info)
	if err != nil {
		log.Println(err)
	}
	return
}

func parseSource(scanner *bufio.Scanner) (source ReportSource) {
	var data string
L:
	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "</source>":
			break L
		default:
			data = data + text + "\n"
		}
	}
	return ReportSource{
		sourceName: "default",
		query:      data,
	}
}

func parseOutput(scanner *bufio.Scanner) (output ReportOutput) {
	openingTag := scanner.Text()
	openingTag = openingTag[1 : len(openingTag)-1]
	tokens := strings.Split(openingTag, " ")
	for _, value := range tokens {
		if value == "csv" || value == "json" || value == "xlsx" {
			output.Format = value
		}
	}
	var data string
L:
	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "</output>":
			break L
		default:
			data = data + text + "\n"
		}
	}
	output.template = data
	if output.Format == "" {
		panic("Not a valid output format")
	}
	if output.template == "" {
		output.template = defaultReports[output.Format]
	}
	return output
}

// MakeReport takes a scanner to a .gore file, reads it and stores in the global Report map
func MakeReport(scanner *bufio.Scanner) (err error) {
	var report Report
	report.Outputs = make(map[string]ReportOutput)

	for scanner.Scan() {
		text := scanner.Text()
		switch {
		case text == "<info>":
			var info ReportInfo
			info, err = parseInfo(scanner)
			if err != nil {
				return
			}
			report.Info = info
		case text == "<source sql>":
			source := parseSource(scanner)
			report.source = source
		case strings.Contains(text, "<output"):
			output := parseOutput(scanner)
			report.Outputs[output.Format] = output
		default:
			log.Fatal("Unrecognized tag " + text)
		}
	}
	Reports[report.Info.Name] = &report
	return
}
