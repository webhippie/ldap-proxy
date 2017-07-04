package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/webhippie/ldap-proxy/assets"
	"github.com/webhippie/ldap-proxy/config"
	"github.com/webhippie/ldap-proxy/router/middleware/header"
	"github.com/webhippie/ldap-proxy/router/middleware/logger"
	"github.com/webhippie/ldap-proxy/router/middleware/recovery"
	"github.com/webhippie/ldap-proxy/templates"
)

// Load initializes the routing of the application.
func Load(middleware ...gin.HandlerFunc) http.Handler {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.SetHTMLTemplate(
		templates.Load(),
	)

	e.Use(middleware...)
	e.Use(logger.SetLogger())
	e.Use(recovery.SetRecovery())
	e.Use(header.SetCache())
	e.Use(header.SetOptions())
	e.Use(header.SetSecure())
	e.Use(header.SetVersion())

	e.StaticFS(
		"/ldap-proxy/assets",
		assets.Load(),
	)

	root := e.Group("/ldap-proxy")
	{
		root.GET("/ping", ping)
	}

	e.NoRoute(func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{},
		)
	})

	if config.Server.Pprof {
		pprof.Register(
			e,
			&pprof.Options{
				RoutePrefix: "/ldap-proxy/debug/pprof",
			},
		)
	}

	return e
}

func ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "pong",
		},
	)
}
