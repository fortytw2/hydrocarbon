package discollect

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"mime"
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
	hash := hex.EncodeToString(h.Sum(nil))

	fName := hash + "." + strings.Split(fileName, ".")[1]

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

	return fmt.Sprintf("%s?file=%s", lf.staticPath, fName), nil
}

// ServeHTTP implements hydrocarbon.ErrorHandler
func (lf *LocalFS) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	fName := r.URL.Query().Get("file")

	ct, _, err := mime.ParseMediaType(fName)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", ct)

	f, err := os.Open(lf.rootPath + "/" + fName)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	return err
}
