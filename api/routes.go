package api

import (
	"net/http"

	"github.com/fortytw2/kiasu"
	"github.com/fortytw2/kiasu/api/middlewarez"
	"github.com/go-kit/kit/log"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

// Routes returns all routes for the HTTP/JSON interface :)
func Routes(l log.Logger, m kiasu.Mailer, ds kiasu.Store) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middlewarez.Logger(l))

	r.Mount("/api/v1/users", userRoutes(l, m, ds))

	return r
}

func userRoutes(l log.Logger, m kiasu.Mailer, ds kiasu.Store) http.Handler {
	users := chi.NewRouter()

	users.Get("/", Authenticate(ds, UserProfile(l, ds)).ServeHTTP)
	users.Delete("/", Authenticate(ds, DeactivateUser(l, ds)).ServeHTTP)

	users.Get("/sessions", Authenticate(ds, UserSessions(l, ds)).ServeHTTP)
	users.Delete("/sessions", Authenticate(ds, Logout(l, ds)).ServeHTTP)

	users.Get("/confirm", ConfirmToken(l, ds).ServeHTTP)
	users.Post("/login", Login(l, ds).ServeHTTP)
	users.Post("/new", RegisterUser(l, m, ds).ServeHTTP)

	return users
}
