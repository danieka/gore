package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/danieka/gore/internal/server"

	"github.com/danieka/gore/internal/interactiveserver"
	"github.com/danieka/gore/internal/reports"
	"github.com/danieka/gore/internal/sources"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// Config is the config for the file
type Config struct {
	Sources map[string]sources.SourceConfig
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var watcher *fsnotify.Watcher

var isGoreFile = regexp.MustCompile(`.gore$`)

var watchMode bool
var interactiveMode bool

func walkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		check(err)
		if !info.IsDir() {
			if isGoreFile.MatchString(path) {
				files = append(files, path)
				if watchMode {
					fmt.Println("adding to watcher", path)
					watcher.Add(path)
				}
			}
		}
		return nil
	})
	return files, err
}

func startWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Reloading report: ", event.Name)
					err := loadReport(event.Name)
					check(err)
					if interactiveMode {
						interactiveserver.TriggerReload()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}

func loadReport(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return reports.MakeReport(scanner)
}

func init() {
	flag.BoolVar(&watchMode, "w", false, "Watch .gore files for changes and reload reports on changes")
	flag.BoolVar(&interactiveMode, "i", false, "Run server in interactive mode (suitable for development)")
}

func main() {
	flag.Parse()

	if watchMode {
		startWatcher()
		defer watcher.Close()
	}

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to open config file: %s", err.Error())
	}

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
		loadReport(path)
		check(err)
	}

	if interactiveMode {
		err = interactiveserver.Start()
		check(err)
	} else {
		err = server.Start()
		check(err)
	}
}
