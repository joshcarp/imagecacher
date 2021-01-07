package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/afero"
)

var (
	input  = flag.String("input", ".", "input directory to replace links in")
	r      = flag.String("regex", `!\[.*?\]\((?P<url>.*?)\)`, "regex link to replace")
	output = flag.String("output", ".", "output directory to put images")
	prefix = flag.String("prefix", ".", "prefix to prefix links with")
)

func main() {
	type replacer struct {
		old, new string
		ext      string
	}
	flag.Parse()
	fs := afero.NewOsFs()
	regex := regexp.MustCompile(*r)
	prereplace := make(map[string]string)
	files := make(map[string][]byte)
	newSvgFiles := make(map[string][]byte)
	replace := make(map[string]replacer)

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
			defer func() { done++ }()
			foundUrl = strings.ReplaceAll(foundUrl, "https://", "http://")
			resp, err := RetryHTTPRequest(foundUrl)
			if err != nil {
				return
			}
			contentType := resp.Header.Get("Content-Type")
			repl, _ := ioutil.ReadAll(resp.Body)
			mimeType, _ := mime.ExtensionsByType(contentType)
			var m string
			for _, e := range mimeType {
				m = e
				break
			}
			replaceChan <- replacer{
				old: src,
				new: string(repl),
				ext: m,
			}
		}(src, foundUrl)
	}

	for {
		select {
		case a, ok := <-replaceChan:
			if !ok {
				goto cont
			}
			replace[a.old] = a
		default:
			if done == len(prereplace) {
				close(replaceChan)
			}
		}
	}
cont:

	for filename := range files {
		for src, newurl := range replace {
			svgname := HashName(src) + newurl.ext
			files[filename] = bytes.ReplaceAll(files[filename], []byte(src), []byte(path.Join(*prefix, svgname)))
			newSvgFiles[svgname] = []byte(newurl.new)
		}
	}

	for filename := range files {
		if err := afero.WriteFile(fs, filename, files[filename], os.ModePerm); err != nil {
			panic(err)
		}
	}
	for filename := range newSvgFiles {
		if err := afero.WriteFile(fs, path.Join(*output, filename), newSvgFiles[filename], os.ModePerm); err != nil {
			panic(err)
		}
	}
	bar.Finish()
}

func HashName(s string) string {
	hash := md5.Sum([]byte(s))
	src := hash[:]
	hexPassword := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(hexPassword, src)
	return string(hexPassword)
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
