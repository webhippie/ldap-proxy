package handler

import (
	"net/http"

	"github.com/unrolled/render"
)

// Proxy redirects to login or proxies the requests.
func Proxy(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// TODO: Add the proxy logic for authenticated users
		http.Redirect(w, req, "/ldap-proxy/login", 301)
	}
}
