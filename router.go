package hydrocarbon

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/fortytw2/hydrocarbon/public"
)

//go:generate bash -c "cd ui && yarn run build-dist"
//go:generate bash -c "go-bindata -pkg public -mode 0644 -modtime 499137600 -o public/assets_generated.go ui/build/..."

func NewRouter(ua *UserAPI) *http.ServeMux {
	m := http.NewServeMux()

	m.Handle("/", http.FileServer(
		&assetfs.AssetFS{Asset: public.Asset, AssetDir: public.AssetDir, AssetInfo: public.AssetInfo, Prefix: "ui/build/"}))

	// session management
	m.HandleFunc("/api/token/request", ua.RequestToken)
	m.HandleFunc("/api/token/activate", ua.Activate)
	m.HandleFunc("/api/key/deactivate", ua.Deactivate)

	return m
}
