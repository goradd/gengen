package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/goradd/gofile/pkg/sys"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var modules map[string]string

func main() {
	var config string
	var outFile string
	var err error

	flag.StringVar(&config, "c", "", "A required config file that will be used to provide the *dot* context to the template.")
	flag.StringVar(&outFile, "o", "", "Output file. If not specified, output will be sent to stdout.")
	flag.Parse() // regular run of program

	if config == "" {
		log.Fatal("you must specify a config file with the -c option.")
	}

	modules, err = sys.ModulePaths()
	if err != nil {
		log.Fatal(err)
	}

	config = getRealPath(config)

	data, err := ioutil.ReadFile(config)
	if err != nil {
		panic(err)
	}

	var dot interface{}

	idx := bytes.IndexRune(data, '{')
	if idx < 0 {
		panic ("The configuration file must contain a json object that starts with an open bracket.")
	}
	err = json.Unmarshal(data[idx:], &dot)
	if err != nil {
		panic(err)
	}

	switch flag.NArg() {
	case 0:
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {panic(err)}
		break
	case 1:
		data, err = ioutil.ReadFile(getRealPath(flag.Arg(0)))
		if err != nil {panic(err)}
		break
	default:
		log.Fatal("input must be from stdin or a single file")
	}

	var tmpl *template.Template
	tmpl,err = template.New("temp").Parse(string(data))
	if err != nil {log.Fatal(err)}

	if outFile == "" {
		err = tmpl.Execute(os.Stdout, dot)
		if err != nil {log.Fatal(err)}
	} else {
		if file,err := os.Create(getRealPath(outFile)); err != nil {
			panic(err)
		} else {
			defer file.Close()
			err = tmpl.Execute(file, dot)
			if err != nil {panic(err)}
		}
	}
}

func getRealPath(path string) string {
	var err error
	path = os.ExpandEnv(path)
	path,err = sys.GetModulePath(path, modules)
	if err != nil {
		log.Fatal(err)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return path
}
