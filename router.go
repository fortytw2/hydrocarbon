package hydrocarbon

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/fortytw2/hydrocarbon/public"
	"github.com/julienschmidt/httprouter"
)

//go:generate bash -c "pushd ui && preact build --service-worker false --no-prerender && popd"
//go:generate bash -c "go-bindata -pkg public -mode 0644 -modtime 499137600 -o public/assets_generated.go ui/build/..."

// ErrorHandler wraps up common error handling patterns for http routers
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (eh ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := eh(w, r)
	if err != nil {
		writeErr(w, err)
	}
}

func limitDecoder(r *http.Request, x interface{}) error {
	return json.NewDecoder(io.LimitReader(r.Body, 1024*8)).Decode(x)
}

var (
	statusOK    = "success"
	statusError = "error"
)

// writeSuccess is a helper for writing the same format of JSON for every reply
func writeSuccess(w http.ResponseWriter, x interface{}) error {
	var s = struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data,omitempty"`
	}{
		statusOK,
		x,
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(s)
}

// writeErr is the only way to write an error
func writeErr(w http.ResponseWriter, uErr error) {
	var s = struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{
		statusError,
		uErr.Error(),
	}
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// NewRouter configures a new http.Handler that serves hydrocarbon
func NewRouter(ua *UserAPI, fa *FeedAPI, domain string) http.Handler {
	m := httprouter.New()

	fs := http.FileServer(
		&assetfs.AssetFS{
			Asset:     public.Asset,
			AssetDir:  public.AssetDir,
			AssetInfo: public.AssetInfo,
			Prefix:    "ui/build/",
		})

	m.Handler("GET", "/static/*file", http.StripPrefix("/static/", fs))

	// serve the single page app for every other route, it has a 404 page builtin
	m.NotFound = ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return nil
		}

		buf := public.MustAsset("ui/build/index.html")
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write(buf)
		return err
	})

	routes := map[string]ErrorHandler{
		// login tokens
		"/v1/token/create": ua.RequestToken,

		// payment managemnet
		"/v1/payment/create": ua.CreatePayment,

		// api keys
		"/v1/key/create": ua.Activate,
		"/v1/key/delete": ua.Deactivate,
		"/v1/key/list":   ua.ListSessions,

		// feed management
		"/v1/feed/create": fa.AddFeed,
		"/v1/feed/delete": fa.RemoveFeed,

		// list all feeds for a folder
		"/v1/feed/list": fa.GetFeedsForFolder,

		// folder management
		"/v1/folder/create": fa.AddFolder,
		// list all folders
		"/v1/folder/list": fa.GetFolders,

		// list all posts in a feed
		"/v1/post/list": fa.GetFeed,
	}

	for route, handler := range routes {
		m.Handler(http.MethodPost, route, handler)
	}

	if httpsOnly(domain) {
		return redirectHTTPS(m)
	}

	return m
}

func httpsOnly(domain string) bool {
	u, err := url.Parse(domain)
	if err != nil {
		panic(err)
	}
	return u.Scheme == "https"
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
