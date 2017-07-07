package handler

import (
	"net/http"

	"github.com/unrolled/render"
	"github.com/webhippie/ldap-proxy/pkg/config"
)

// Login displays the login form for authentication.
func Login(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(
			w,
			http.StatusOK,
			"login",
			map[string]string{
				"Title": config.Server.Title,
			},
		)
	}
}
