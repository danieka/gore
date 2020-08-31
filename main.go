package main

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/danieka/gore/internal/interactiveserver"
	"github.com/danieka/gore/internal/reports"
	"github.com/danieka/gore/internal/sources"

	"gopkg.in/yaml.v2"
)

// Config is the config for the file
type Config struct {
	Sources map[string]sources.SourceConfig
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data, err := ioutil.ReadFile("config.yaml")
	check(err)

	config := Config{}

	err = yaml.Unmarshal([]byte(data), &config)
	check(err)

	for k, v := range config.Sources {
		sources.MakeSQLSource(k, v)
	}

	file, err := os.Open("test.gore")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	err = reports.MakeReport(scanner)
	check(err)

	err = interactiveserver.Start()
	check(err)
}
