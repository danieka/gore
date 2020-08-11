package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// InputConfig defines a source
type InputConfig struct {
	Type     string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

// Config is the config for the file
type Config struct {
	Inputs map[string]InputConfig
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var inputs map[string]Input = make(map[string]Input)

func main() {
	data, err := ioutil.ReadFile("config.yaml")
	check(err)

	config := Config{}

	err = yaml.Unmarshal([]byte(data), &config)
	check(err)

	for k, v := range config.Inputs {
		inputs[k] = MakeSQLInput(v)
	}

	file, err := os.Open("test.rpt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	_, err = MakeReport(scanner)
	check(err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	fmt.Println("Starting server")
	err = http.ListenAndServe(":16772", nil)
	check(err)
}
