package handler

import (
	"net/http"
	"path"

	"github.com/go-kit/kit/log"
	"github.com/webhippie/ldap-proxy/pkg/assets"
	"github.com/webhippie/ldap-proxy/pkg/config"
)

// Static handles all requests to static assets.
func Static(logger log.Logger) http.Handler {
	return http.StripPrefix(
		path.Join(
			config.Server.Root,
			"assets",
		),
		http.FileServer(
			assets.Load(logger),
		),
	)
}
