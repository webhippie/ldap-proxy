package handler

import (
	"net/http"
	"path"

	"github.com/go-kit/kit/log"
	"github.com/webhippie/ldap-proxy/pkg/config"
)

// Proxy redirects to login or proxies the requests.
func Proxy(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO: Add the proxy logic for authenticated users

		http.Redirect(
			w,
			r,
			path.Join(
				config.Server.Root,
				"login",
			),
			http.StatusMovedPermanently,
		)
	}
}
