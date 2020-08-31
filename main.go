package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

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

var isGoreFile = regexp.MustCompile(`.gore$`)

func walkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		check(err)
		if !info.IsDir() {
			if isGoreFile.MatchString(path) {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
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

	wd, err := os.Getwd()
	check(err)

	files, err := walkDir(wd)

	for _, path := range files {
		file, err := os.Open(path)
		check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		err = reports.MakeReport(scanner)
		check(err)
	}

	err = interactiveserver.Start()
	check(err)
}
