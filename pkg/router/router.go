package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/unrolled/render"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/handler"
	"github.com/webhippie/ldap-proxy/pkg/router/middleware/header"
	"github.com/webhippie/ldap-proxy/pkg/router/middleware/prometheus"
	"github.com/webhippie/ldap-proxy/pkg/templates"
)

// Load initializes the routing of the application.
func Load() http.Handler {
	r := render.New(render.Options{
		Asset:         templates.File,
		AssetNames:    templates.Names,
		IsDevelopment: config.Debug,
		Layout:        "layout",
	})

	mux := chi.NewRouter()

	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.NotFound(handler.Proxy(r))

	mux.Route("/ldap-proxy", func(root chi.Router) {
		root.Handle("/assets/*", handler.Static())
		root.Get("/ping", handler.Ping(r))
		root.Get("/login", handler.Login(r))
		root.Post("/login", handler.Auth(r))

		if config.Server.Prometheus {
			root.Get("/metrics", prometheus.Handler())
		}

		if config.Server.Pprof {
			root.Mount("/debug", middleware.Profiler())
		}
	})

	return mux
}
