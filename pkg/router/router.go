package router

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/fail"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/handler"
	"github.com/webhippie/ldap-proxy/pkg/middleware/header"
)

// Load initializes the routing of the application.
func Load(cfg *config.Config, proxy http.Handler) http.Handler {
	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.URLHandler("path"))
	mux.Use(hlog.MethodHandler("method"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(middleware.RealIP)

	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.NotFound(handler.Proxy(cfg, proxy))

	mux.Route(cfg.Server.Root, func(root chi.Router) {
		root.Get("/login", handler.Login(cfg))
		root.Post("/login", handler.Auth(cfg))

		root.Handle("/assets/*", handler.Static(cfg))
	})

	return mux
}

// Status initializes the routing of metrics and healtchecks.
func Status(cfg *config.Config) http.Handler {
	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.URLHandler("path"))
	mux.Use(hlog.MethodHandler("method"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(middleware.RealIP)

	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.Route("/", func(root chi.Router) {
		root.Mount("/metrics", promhttp.Handler())

		root.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})

		root.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})
	})

	return mux
}

// Redirect handles HTTP to HTTPS redirecting.
func Redirect(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsed, err := url.Parse(
			cfg.Server.Host,
		)

		if err != nil {
			log.Info().
				Err(err).
				Msg("failed to parse host")

			fail.ErrorPlain(w, fail.Cause(err).Unexpected())
		}

		target := strings.Join(
			[]string{
				"https://",
				parsed.Host,
				r.URL.Path,
			},
			"",
		)

		if len(r.URL.RawQuery) > 0 {
			target += "?" + r.URL.RawQuery
		}

		log.Info().
			Str("source", r.URL.String()).
			Str("target", target).
			Msg("redirecting to https")

		http.Redirect(
			w,
			r,
			target,
			http.StatusPermanentRedirect,
		)
	})
}
