package handler

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/webhippie/fail"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/templates"
)

// Login displays the login form for authentication.
func Login(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := map[string]string{
			"Title": cfg.Proxy.Title,
			"Root":  cfg.Server.Root,
			"Error": "",
		}

		if err := templates.Load(cfg).ExecuteTemplate(w, "login.tmpl", vars); err != nil {
			log.Warn().
				Err(err).
				Msg("failed to process login template")

			fail.ErrorPlain(w, fail.Cause(err).Unexpected())
			return
		}
	}
}
