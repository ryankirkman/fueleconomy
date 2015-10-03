package workers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	FILE_URL   = "https://www.fueleconomy.gov/feg/epadata/%s.xml.zip"
	UNZIP_PATH = "/tmp/%s"
	ZIP_PATH   = "/tmp/%s.xml.zip"
)

func DownloadXml(name string) (string, error) {
	zipUrl := fmt.Sprintf(FILE_URL, name)
	zipPath := fmt.Sprintf(ZIP_PATH, name)

	unzipPath := fmt.Sprintf(UNZIP_PATH, name)

	err := downloadFile(zipUrl, zipPath)
	if err != nil {
		return "", err
	}
	defer os.Remove(zipPath)

	err = unzip(zipPath, unzipPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s.xml", unzipPath, name), nil
}

func downloadFile(url string, fpath string) error {
	out, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			panic(err)
		}
	}()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileHeader.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
