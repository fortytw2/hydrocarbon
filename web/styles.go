package web

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/fortytw2/httpkit"
)

var styleHash string

func init() {
	buf, _ := Asset("dist/hydrocarbon.min.css")
	hasher := md5.New()
	_, _ = hasher.Write(buf)
	styleHash = hex.EncodeToString(hasher.Sum(nil))
}

//go:generate bash -c "lessc --clean-css less/base.less dist/hydrocarbon.min.css"
//go:generate go-bindata -pkg web -o styles_generated.go -ignore .gitkeep dist/

// Stylesheet writes the stylesheet
func Stylesheet(w http.ResponseWriter, r *http.Request) error {
	buf, err := Asset("dist/hydrocarbon.min.css")
	if err != nil {
		return httpkit.Wrap(err, 404)
	}
	w.Header().Set("Content-Type", "text/css")

	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Set("ETag", styleHash)

	_, err = w.Write(buf)
	if err != nil {
		panic(err)
	}

	return nil
}
