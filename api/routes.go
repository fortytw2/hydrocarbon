package api

import (
	"net/http"

	"github.com/fortytw2/kiasu"
	"github.com/go-kit/kit/log"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

// Routes returns all routes for the HTTP/JSON interface :)
func Routes(l log.Logger, m kiasu.Mailer, ds kiasu.Store) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/v1/users", userRoutes(l, m, ds))

	return r
}

func userRoutes(l log.Logger, m kiasu.Mailer, ds kiasu.Store) http.Handler {
	users := chi.NewRouter()

	users.Get("/", UserProfile(l, ds))
	users.Delete("/", DeactivateUser(l, ds))

	users.Get("/sessions", UserSessions(l, ds))
	users.Delete("/sessions", Logout(l, ds))

	users.Get("/confirm", ConfirmToken(l, m, ds))
	users.Post("/login", Login(l, m, ds))
	users.Post("/new", RegisterUser(l, m, ds))

	return users
}
