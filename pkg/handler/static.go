package handler

import (
	"net/http"
	"path"

	"github.com/webhippie/ldap-proxy/pkg/assets"
	"github.com/webhippie/ldap-proxy/pkg/config"
)

// Static handles all requests to static assets.
func Static(cfg *config.Config) http.Handler {
	return http.StripPrefix(
		path.Join(
			cfg.Server.Root,
			"assets",
		),
		http.FileServer(
			assets.Load(cfg),
		),
	)
}
