package gengen

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"go/build"
	"io/ioutil"
	"encoding/json"
	"text/template"
	"bytes"
)


func main() {
	var config string
	var outFile string

	flag.StringVar(&config, "c", "", "A required config file that will be used to provide the *dot* context to the template.")
	flag.StringVar(&outFile, "o", "", "Output file. If not specified, output will be sent to stdout.")
	flag.Parse() // regular run of program

	if config == "" {
		panic("you must specify a config file use the -c option.")
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
		panic("input must be from stdin or a single file")
		os.Exit(1)
	}

	var tmpl *template.Template
	tmpl,err = template.New("temp").Parse(string(data))
	if err != nil {panic(err)}

	if outFile == "" {
		tmpl.Execute(os.Stdout, dot)
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
	path = filepath.FromSlash(path)
	if strings.Index(path, "GOPATH") == 0 {
		path = goPath() + path[6:]
	}

	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return path
}

func goPath() string {
	var path string
	goPaths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(goPaths) == 0 {
		path = build.Default.GOPATH
	} else if goPaths[0] == "" {
		path = build.Default.GOPATH
	} else {
		path = goPaths[0]
	}

	// clean path so it does not end with a path separator
	if path[len(path)-1] == os.PathSeparator {
		path = path[:len(path)-1]
	}

	// If the GOPATH is empty, then see if the current executable looks like it is in a project
	if path == "" {
		if path2, err := os.Executable(); err == nil {
			path2 = filepath.Join(filepath.Dir(filepath.Dir(path2)), "src")
			dstInfo, err := os.Stat(path)
			if err == nil && dstInfo.IsDir() {
				path = path2
			}
		}
	}

	path,_ = filepath.Abs(path)

	// TODO: GoPath may go away, so we might need to use another way to search for the current go project structure
	return path
}

