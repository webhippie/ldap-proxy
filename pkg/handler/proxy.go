package handler

import (
	"net/http"
	"path"

	"github.com/webhippie/ldap-proxy/pkg/config"
)

// Proxy redirects to login or proxies the requests.
func Proxy(cfg *config.Config, proxy http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add(cfg.Proxy.UserHeader, "myuser")

		// TODO: check if user is authenticated

		http.Redirect(
			w,
			r,
			path.Join(
				cfg.Server.Root,
				"login",
			),
			http.StatusMovedPermanently,
		)

		// proxy.ServeHTTP(w, r)
	}
}
