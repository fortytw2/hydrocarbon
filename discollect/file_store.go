package discollect

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// A FileStore shoves files somewhere and returns a link at which they can be
// retrieved
type FileStore interface {
	Put(fileName string, contents []byte) (string, error)
}

// NewStubFS is used only for testing and doesn't actually do anything
func NewStubFS() *StubFS {
	return &StubFS{
		URL: "https://stubfotos.com/",
	}
}

type StubFS struct {
	URL string
}

func (sf *StubFS) Put(fileName string, contents []byte) (string, error) {
	return sf.URL + fileName, nil
}

// NewLocalFS creates a LocalFS set up to save files to path and serve from
// staticPath
func NewLocalFS(path, staticPath string) (*LocalFS, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(abs, 0644)
	if err != nil {
		return nil, err
	}

	return &LocalFS{
		rootPath:   abs,
		staticPath: staticPath,
	}, nil
}

// LocalFS is both a FileStore implementation backed by the filesystem, but also
// a http.Handler that will serve the images back up
type LocalFS struct {
	rootPath   string
	staticPath string
}

// Put writes the file to disk after hashing it
func (lf *LocalFS) Put(fileName string, contents []byte) (string, error) {
	h := sha1.New()
	_, err := h.Write(contents)
	if err != nil {
		return "", err
	}
	hash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	contentType := http.DetectContentType(contents)

	// only hash the file for now
	fName := hash
	switch contentType {
	case "image/png":
		fName += ".png"
	case "image/jpeg":
		fName += ".jpeg"
	case "image/gif":
		fName += ".gif"
	case "image/webp":
		fName += ".webp"
	default:
		return "", fmt.Errorf("unsupported image type: %s", contentType)
	}

	f, err := os.OpenFile(lf.rootPath+"/"+fName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}

	_, err = f.Write(contents)
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", lf.staticPath, fName), nil
}

// ServeHTTP implements hydrocarbon.ErrorHandler
func (lf *LocalFS) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	splitPath := strings.Split(r.URL.Path, "/")
	fName := splitPath[len(splitPath)-1]
	filePath := lf.rootPath + "/" + fName

	if etag := r.Header.Get("ETag"); etag != "" {
		_, err := os.Stat(filePath)
		if err != nil {
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}

	suffix := strings.Split(fName, ".")[len(strings.Split(fName, "."))-1]
	w.Header().Set("Content-Type", "image/"+suffix)
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	w.Header().Set("Cache-Control", "public, immutable, max-age=31536000")
	w.Header().Set("ETag", fName)

	_, err = w.Write(buf)
	return err
}
