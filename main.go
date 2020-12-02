package main

import (
	"bytes"
	"github.com/alecthomas/kingpin"
	"github.com/cheggaaa/pb/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/afero"
	"github.com/vincent-petithory/dataurl"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	input = kingpin.Arg("input", "input directory to replace links in").Required().String()
)

func main() {
	kingpin.Parse()
	fs := afero.NewOsFs()
	regex := regexp.MustCompile(`!\[.*?\]\((?P<url>.*?)\)`)
	prereplace := make(map[string]string)
	files := make(map[string][]byte)
	replace := make(map[string]string)

	type replacer struct {
		old, new string
	}

	replaceChan := make(chan replacer)
	afero.Walk(fs, *input, func(path string, info os.FileInfo, err error) error {
		b, err := afero.ReadFile(fs, path)
		if err != nil {
			return nil
		}
		for _, match := range regex.FindAllStringSubmatch(string(b), -1) {
			files[path] = b
			for i, name := range regex.SubexpNames() {
				if name == "url" && match[i] != "" {
					foundUrl := strings.TrimLeft(strings.TrimRight(match[i], `"`), `"`)
					if u, err := url.ParseRequestURI(foundUrl); err != nil || u.Host == "" {
						continue
					}
					prereplace[foundUrl] = foundUrl
				}
			}
		}
		return nil
	})
	bar := pb.StartNew(len(prereplace))
	done := 0
	for src, foundUrl := range prereplace {

		go func(src, foundUrl string) {
			defer bar.Increment()
			defer func(){done++}()
			foundUrl = strings.ReplaceAll(foundUrl, "https://", "http://")
			resp, err := RetryHTTPRequest(foundUrl)
			if err != nil {
				return
			}
			contentType := resp.Header.Get("Content-Type")
			repl, _ := ioutil.ReadAll(resp.Body)
			durl := dataurl.New(repl, contentType)
			replaceChan <- replacer{
				old: src,
				new: durl.String(),
			}
		}(src, foundUrl)
	}

	for {
		select {
		case a, ok := <-replaceChan:
			if !ok {
				goto cont
			}
			replace[a.old] = a.new
		default:
			if done == len(prereplace){
				close(replaceChan)
			}
		}
	}
	cont:

	for filename := range files {
		for src, newurl := range replace {
			files[filename] = bytes.ReplaceAll(files[filename], []byte(src), []byte(newurl))
		}
	}

	for filename := range files {
		if err := afero.WriteFile(fs, filename, files[filename], os.ModePerm); err != nil {
			panic(err)
		}
	}
	bar.Finish()
}

// RetryHTTPRequest retries the given request
func RetryHTTPRequest(url string) (*http.Response, error) {
	client := retryablehttp.NewClient()
	client.Logger = nil
	client.RetryMax = 100
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}