package common

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// extract tar.gz files
func Ungzip(source, target string) error {

	// create a file descriptor for the tar.gz file (in RO mode)
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// create a new gzip reader
	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	// create a new tar reader
	tarReader := tar.NewReader(archive)

	for {

		// advances to the next entry in the tar archive
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// join the path of the last file found from tar with "target"
		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		// create a fd for the extracted file in write mode
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		// make an io copy with the extracted file from tar and the "file" fd
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return err

}

func FileSize(url string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Header.Get("Content-Length"), nil
}
