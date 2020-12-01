package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/spf13/afero"
	"github.com/vincent-petithory/dataurl"
)

var (
	input = kingpin.Arg("input", "input directory to replace links in").Required().String()
)

func main() {
	kingpin.Parse()
	fs := afero.NewOsFs()
	regex := regexp.MustCompile(`(?:<img.*?src\s*=\s*)(?P<url>".*?")`)
	replace := make(map[string]string)
	files := make(map[string][]byte)

	afero.Walk(fs, *input, func(path string, info os.FileInfo, err error) error {
		b, err := afero.ReadFile(fs, path)
		if err != nil {
			return nil
		}
		for _, match := range regex.FindAllStringSubmatch(string(b), -1) {
			files[path] = b
			for i, name := range regex.SubexpNames() {
				if name == "url" && match[i] != "" {
					replace[match[i]] = match[i]
				}
			}
		}
		return nil
	})

	for src, url := range replace {
		url = strings.TrimLeft(strings.TrimRight(url, `"`), `"`)
		resp, err := http.Get(url)
		if err == nil {
			contentType := resp.Header.Get("Content-Type")
			repl, _ := ioutil.ReadAll(resp.Body)
			durl := dataurl.New(repl, contentType)
			replace[src] = `"` + durl.String() + `"`
		}
	}

	for src, url := range replace {
		for filename, content := range files {
			files[filename] = bytes.ReplaceAll(content, []byte(src), []byte(url))
		}
	}
	for filename, content := range files {
		if err := afero.WriteFile(fs, filename, content, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
