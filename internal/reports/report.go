package reports

import (
	"bufio"
	"fmt"
	"log"

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

// Report is the main struct for a report
type Report struct {
	Info   ReportInfo
	source ReportSource
}

// Execute the report and return the output
func (r *Report) Execute() (s string, err error) {
	source := sources.Sources[r.source.sourceName]
	if source == nil {
		return "", fmt.Errorf("Unable to find source %s", r.source.sourceName)
	}
	rows, err := sources.Sources[r.source.sourceName].Execute(r.source.query)
	return fmt.Sprintf("%v", rows), err
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

// MakeReport takes a scanner to a .rpt file, reads it and stores in the global Report map
func MakeReport(scanner *bufio.Scanner) (err error) {
	var report Report
	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "<info>":
			var info ReportInfo
			info, err = parseInfo(scanner)
			if err != nil {
				return
			}
			report.Info = info
		case "<source sql>":
			source := parseSource(scanner)
			report.source = source
		default:
			log.Fatal("Unrecognized tag " + text)
		}
	}
	Reports[report.Info.Name] = &report
	return
}
