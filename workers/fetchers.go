package workers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Fetcher interface {
	Fetch(string) ([]byte, error)
}

type RestFetcher struct{}

func (r RestFetcher) Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type FileFetcher struct{}

func (v FileFetcher) Fetch(fname string) ([]byte, error) {
	xmlFilePath, err := DownloadXml(fname)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(filepath.Dir(xmlFilePath))
	xmlFile, err := os.Open(xmlFilePath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	return b, nil
}
