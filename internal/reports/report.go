package reports

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/danieka/gore/internal/sources"

	"gopkg.in/yaml.v2"
)

// ReportInfo contains all metadata about a report
type ReportInfo struct {
	Name string
}

// ReportSource is the data source for the report
type ReportSource struct {
	sourceName string
	query      string
}

// ReportOutput contains info on the output
type ReportOutput struct {
	format   string
	template string
}

// Report is the main struct for a report
type Report struct {
	Info    ReportInfo
	source  ReportSource
	outputs map[string]ReportOutput
}

// Execute the report and return the output
func (r *Report) Execute(format string) (s string, err error) {
	source := sources.Sources[r.source.sourceName]
	if source == nil {
		return "", fmt.Errorf("Unable to find source %s", r.source.sourceName)
	}
	cols, rows, err := sources.Sources[r.source.sourceName].Execute(r.source.query)

	output, ok := r.outputs[format]
	if !ok {
		return "", fmt.Errorf("Unable to find output format %s", format)
	}

	tmpl, err := template.New("").Parse(output.template)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "", map[string]interface{}{
		"Rows": rows,
		"Cols": cols,
	})
	if err != nil {
		log.Println(err)
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
		if value == "csv" || value == "json" {
			output.format = value
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
	if output.format == "" {
		panic("No valid output format")
	}
	if output.template == "" {
		output.template = defaultReports[output.format]
	}
	return output
}

// MakeReport takes a scanner to a .gore file, reads it and stores in the global Report map
func MakeReport(scanner *bufio.Scanner) (err error) {
	var report Report
	report.outputs = make(map[string]ReportOutput)

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
			report.outputs[output.format] = output
		default:
			log.Fatal("Unrecognized tag " + text)
		}
	}
	Reports[report.Info.Name] = &report
	return
}
