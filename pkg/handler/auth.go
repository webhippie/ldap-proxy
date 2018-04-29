package handler

import (
	"net/http"
	"strings"

	"github.com/codehack/fail"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/templates"
)

// Auth handles the authentication itself against LDAP.
func Auth(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "auth")

	return func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			level.Warn(logger).Log(
				"msg", "failed to parse login form",
				"err", err,
			)

			fail.Error(w, fail.Cause(err).BadRequest("failed to parse login form"))
			return
		}

		conn, err := ldap.Dial("tcp", config.LDAP.Addr)

		if err != nil {
			level.Warn(logger).Log(
				"msg", "failed to connect to ldap",
				"err", err,
			)

			fail.Error(w, fail.Cause(err).BadRequest("failed to connect to ldap"))
		}

		defer conn.Close()

		err := conn.Bind(
			config.LDAP.BindUsername,
			config.LDAP.BindPassword,
		)

		if err != nil {
			level.Warn(logger).Log(
				"msg", "failed to bind to ldap",
				"err", err,
			)

			fail.Error(w, fail.Cause(err).BadRequest("failed to bind to ldap"))
		}

		result, err := conn.Search(ldap.NewSearchRequest(
			config.LDAP.BaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			strings.Replace(
				config.LDAP.FilterDN,
				"{login}",
				req.PostFormValue("username"),
				-1,
			),
			[]string{"dn", config.LDAP.UserAttr},
			nil,
		))

		if err != nil {
			level.Warn(logger).Log(
				"msg", "failed to find a user",
				"err", err,
			)

			fail.Error(w, fail.Cause(err).BadRequest("failed to find a user"))
			return
		}

		if len(result.Entries) < 1 {
			level.Warn(logger).Log(
				"msg", "user does not exist",
				"user", req.PostFormValue("username"),
				"count", len(result.Entries),
			)

			fail.Error(w, fail.Cause(err).BadRequest("user does not exist"))
			return
		}

		if len(result.Entries) > 1 {
			level.Warn(logger).Log(
				"msg", "too many user matches",
				"user", req.PostFormValue("username"),
				"count", len(result.Entries),
			)

			fail.Error(w, fail.Cause(err).BadRequest("too many user matches"))
			return
		}

		if err := conn.Bind(result.Entries[0].DN, req.PostFormValue("password")); err != nil {
			level.Warn(logger).Log(
				"msg", "failed to authenticate user",
				"err", err,
			)

			fail.Error(w, fail.Cause(err).BadRequest("failed to authenticate user"))
			return
		}

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
