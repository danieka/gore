package main

import (
	"bufio"
	"log"

	"gopkg.in/yaml.v2"
)

// ReportInfo contains all metadata about a report
type ReportInfo struct {
	Name string
}

// Report is the main struct for a report
type Report struct {
	Info ReportInfo
}

func parseInfo(scanner *bufio.Scanner) (info ReportInfo, err error) {
	var data string
	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "</info>":
			break
		default:
			data = data + text + "\n"
		}
	}
	err = yaml.Unmarshal([]byte(data), &info)
	check(err)
	return
}

// MakeReport takes a scanner to a .rpt file, reads it and returns the report
func MakeReport(scanner *bufio.Scanner) (report Report, err error) {
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
		default:
			log.Fatal("Unrecognized tag " + text)
		}
	}
	return
}
