package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/handler"
	"github.com/webhippie/ldap-proxy/pkg/router/middleware/header"
	"github.com/webhippie/ldap-proxy/pkg/router/middleware/prometheus"
	"github.com/webhippie/ldap-proxy/pkg/router/middleware/requests"
)

// Load initializes the routing of the application.
func Load(logger log.Logger) http.Handler {
	mux := chi.NewRouter()

	mux.Use(requests.Requests(logger))

	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(middleware.RealIP)

	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.NotFound(handler.Proxy(logger))

	mux.Route(config.Server.Root, func(root chi.Router) {
		if config.Server.Prometheus {
			root.Get("/metrics", prometheus.Handler())
		}

		if config.Server.Pprof {
			root.Mount("/debug", middleware.Profiler())
		}

		root.Get("/healthz", handler.Healthz(logger))
		root.Get("/readyz", handler.Readyz(logger))

		root.Get("/login", handler.Login(logger))
		root.Post("/login", handler.Auth(logger))

		root.Handle("/assets/*", handler.Static(logger))
	})

	return mux
}
