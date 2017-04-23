package hydrocarbon

import (
	"net/http"

	"github.com/bouk/httprouter"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/fortytw2/hydrocarbon/public"
)

//go:generate bash -c "cd ui && yarn run build-dist"
//go:generate bash -c "go-bindata -pkg public -mode 0644 -modtime 499137600 -o public/assets_generated.go ui/build/..."

func NewRouter(ua *UserAPI, https bool) http.Handler {
	m := httprouter.New()

	fs := http.FileServer(
		&assetfs.AssetFS{
			Asset:     public.Asset,
			AssetDir:  public.AssetDir,
			AssetInfo: public.AssetInfo,
			Prefix:    "ui/build/",
		})

	m.Handle("GET", "/static/*file", http.StripPrefix("/static/", fs).ServeHTTP)
	m.Handle("GET", "/favicon.ico", fs.ServeHTTP)

	m.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := public.MustAsset("ui/build/index.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(buf)
	})

	// session management
	m.POST("/api/token/request", ua.RequestToken)
	m.POST("/api/token/activate", ua.Activate)
	m.POST("/api/key/deactivate", ua.Deactivate)

	if https {
		return redirectHTTPS(m)
	}
	return m
}

func redirectHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			r.URL.Scheme = "https"
			http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
