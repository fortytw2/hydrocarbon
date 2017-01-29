package web

import (
	"net"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/httputil"
	"github.com/fortytw2/hydrocarbon/internal/log"
	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/unrolled/secure"
)

const (
	homeURL           = "/"
	loginURL          = "/login"
	logoutURL         = "/logout"
	privacyURL        = "/privacy"
	registerURL       = "/register"
	confirmTokenURL   = "/confirm"
	forgotPasswordURL = "/password_reset"

	feedsURL        = "/feeds"
	newFeedURL      = "/feeds/new"
	reorderFeedsURL = "/reorder_feeds"

	onePostURL    = "/posts/:postID"
	readStatusURL = "/mark_read"

	settingsURL = "/settings"
)

//go:generate ftmpl -targetgo ./templates_generated.go templates/
//go:generate goimports -w ./templates_generated.go

// Routes returns all routes for this application
func Routes(s *hydrocarbon.Store, l log.Logger, db *geoip2.Reader) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(secureHeader(true))
	r.Use(httplog(l, db))
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.DefaultCompress)

	r.Handle("/hydrocarbon.min.css", httputil.ErrorHandler(Stylesheet).Func())

	r.With(authenticate(s, l)).Get(homeURL, httputil.ErrorHandler(renderHome).Func())
	r.With(authenticate(s, l)).Get(privacyURL, httputil.ErrorHandler(renderPrivacy).Func())

	r.With(authenticate(s, l)).Get(registerURL, httputil.ErrorHandler(renderRegister).Func())
	r.Post(registerURL, newUser(s).Func())
	r.With(authenticate(s, l)).Get(confirmTokenURL, confirmUser)
	r.With(authenticate(s, l)).Get(settingsURL, httputil.ErrorHandler(renderSettings).Func())

	r.Get(forgotPasswordURL, httputil.ErrorHandler(renderPasswordReset).Func())
	r.Post(forgotPasswordURL, forgotPassword)

	r.With(authenticate(s, l)).Get(loginURL, httputil.ErrorHandler(renderLogin).Func())
	r.Post(loginURL, httputil.ErrorHandler(newSession(s)).Func())
	r.With(authenticate(s, l)).Get(logoutURL, deleteSession(s).Func())

	r.With(authenticate(s, l)).Get(feedsURL, renderFeed(s).Func())

	r.With(authenticate(s, l)).Get(newFeedURL, renderNewFeed)
	r.With(authenticate(s, l)).Post(newFeedURL, addFeed(s).Func())
	r.With(authenticate(s, l)).Post(reorderFeedsURL, reorderFeeds)
	r.With(authenticate(s, l)).Delete(feedsURL, deleteFeed)

	r.With(authenticate(s, l)).Post("/charge", activateAccount(s, l).Func())
	r.Handle("/stripe_webhooks", stripeWebhookHandler(s, l))

	r.Get(onePostURL, renderPost)
	r.Post(readStatusURL, markRead)

	return r
}

// secureHeader sets security headers
func secureHeader(dev bool) func(http.Handler) http.Handler {
	s := secure.New(secure.Options{
		AllowedHosts:            []string{"hydrocarbon.io"},                      // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		SSLRedirect:             true,                                            // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                           // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "hydrocarbon.io",                                // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"}, // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              315360000,                                       // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                            // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                            // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          false,                                           // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                            // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                    // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option.
		ContentTypeNosniff:      true,                                            // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:        true,                                            // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		// PublicKey:               `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://6e388a801d574fd8d8fbd52f05cf2ae9.report-uri.io/r/default/csp/enforce"`, // PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".

		IsDevelopment: dev, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	return s.Handler
}

func httplog(l log.Logger, db *geoip2.Reader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			m := httpsnoop.CaptureMetrics(next, w, req)
			var remoteAddr string
			if req.Header.Get("x-real-ip") != "" {
				remoteAddr = req.Header.Get("x-real-ip")
			} else if req.Header.Get("x-forwarded-for") != "" {
				remoteAddr = req.Header.Get("x-forwarded-for")
			} else if req.Header.Get("x-client-ip") != "" {
				remoteAddr = req.Header.Get("x-client-ip")
			} else {
				remoteAddr = req.RemoteAddr
			}

			ip := net.ParseIP(remoteAddr)
			c, err := db.Country(ip)
			if err != nil {
				l.Log("msg", "request", "method", req.Method, "url", req.URL.String(), "ip", remoteAddr, "code", m.Code, "time", m.Duration, "bytes", m.Written)
			} else {
				l.Log("msg", "request", "method", req.Method, "url", req.URL.String(), "ip", remoteAddr, "country", c.Country.IsoCode, "code", m.Code, "time", m.Duration, "bytes", m.Written)
			}
		}

		return http.HandlerFunc(fn)
	}
}
