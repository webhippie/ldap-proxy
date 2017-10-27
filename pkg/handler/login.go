package handler

import (
	"net/http"

	"github.com/codehack/fail"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/templates"
)

// Login displays the login form for authentication.
func Login(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "login")

	return func(w http.ResponseWriter, req *http.Request) {
		err := templates.Load(logger).ExecuteTemplate(
			w,
			"login.tmpl",
			map[string]string{
				"Title": config.Server.Title,
				"Root":  config.Server.Root,
			},
		)

		if err != nil {
			level.Warn(logger).Log(
				"msg", "failed to process index template",
				"err", err,
			)

			fail.Error(w, fail.Cause(err).Unexpected())
			return
		}
	}
}
