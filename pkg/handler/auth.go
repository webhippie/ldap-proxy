package handler

import (
	"net/http"
	"path"

	"github.com/webhippie/ldap-proxy/pkg/config"
)

// Auth handles the authentication itself against LDAP.
func Auth(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(
			w,
			r,
			path.Join(
				cfg.Server.Root,
				"login",
			),
			http.StatusMovedPermanently,
		)
	}
}
