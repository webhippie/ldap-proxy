package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/webhippie/ldap-proxy/config"
	"github.com/webhippie/ldap-proxy/router"
	"gopkg.in/urfave/cli.v2"

	// Import pprof for optional debugging
	_ "net/http/pprof"
)

// Server provides the sub-command to start the server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the integrated server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				Value:       "0.0.0.0:8080",
				Usage:       "Address to bind the server",
				EnvVars:     []string{"LDAP_PROXY_ADDR"},
				Destination: &config.Server.Addr,
			},
			&cli.BoolFlag{
				Name:        "pprof",
				Usage:       "Enable pprof debugging server",
				EnvVars:     []string{"LDAP_PROXY_PPROF"},
				Destination: &config.Server.Pprof,
			},
			&cli.StringFlag{
				Name:        "cert",
				Value:       "",
				Usage:       "Path to SSL cert",
				EnvVars:     []string{"LDAP_PROXY_CERT"},
				Destination: &config.Server.Cert,
			},
			&cli.StringFlag{
				Name:        "key",
				Value:       "",
				Usage:       "Path to SSL key",
				EnvVars:     []string{"LDAP_PROXY_KEY"},
				Destination: &config.Server.Key,
			},
			&cli.StringFlag{
				Name:        "templates",
				Value:       "",
				Usage:       "Path to custom templates",
				EnvVars:     []string{"LDAP_PROXY_TEMPLATES"},
				Destination: &config.Server.Templates,
			},
			&cli.StringFlag{
				Name:        "assets",
				Value:       "",
				Usage:       "Path to custom assets",
				EnvVars:     []string{"LDAP_PROXY_ASSETS"},
				Destination: &config.Server.Assets,
			},
		},
		Action: func(c *cli.Context) error {
			logrus.Infof("Starting on %s", config.Server.Addr)

			var (
				server *http.Server
			)

			if config.Server.Cert != "" && config.Server.Key != "" {
				cert, err := tls.LoadX509KeyPair(
					config.Server.Cert,
					config.Server.Key,
				)

				if err != nil {
					logrus.Fatal("Failed to load SSL certificates. %s", err)
				}

				curves := []tls.CurveID{
					tls.CurveP521,
					tls.CurveP384,
					tls.CurveP256,
				}

				ciphers := []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				}

				cfg := &tls.Config{
					PreferServerCipherSuites: true,
					MinVersion:               tls.VersionTLS12,
					CurvePreferences:         curves,
					CipherSuites:             ciphers,
					Certificates:             []tls.Certificate{cert},
				}

				server = &http.Server{
					Addr:         config.Server.Addr,
					Handler:      router.Load(),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					TLSConfig:    cfg,
				}
			} else {
				server = &http.Server{
					Addr:         config.Server.Addr,
					Handler:      router.Load(),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}
			}

			if err := startServer(server); err != nil {
				logrus.Fatal(err)
			}

			return nil
		},
	}
}
