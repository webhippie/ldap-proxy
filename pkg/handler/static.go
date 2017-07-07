package handler

import (
	"net/http"

	"github.com/webhippie/ldap-proxy/pkg/assets"
)

// Static handles all requests to static assets.
func Static() http.Handler {
	return http.StripPrefix(
		"/ldap-proxy/assets",
		http.FileServer(
			assets.Load(),
		),
	)
}
