package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/router"
	"gopkg.in/urfave/cli.v2"
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
				Value:       false,
				Usage:       "Enable pprof debugging server",
				EnvVars:     []string{"LDAP_PROXY_PPROF"},
				Destination: &config.Server.Pprof,
			},
			&cli.BoolFlag{
				Name:        "prometheus",
				Value:       false,
				Usage:       "Enable prometheus exporter",
				EnvVars:     []string{"LDAP_PROXY_PROMETHEUS"},
				Destination: &config.Server.Prometheus,
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
			&cli.StringFlag{
				Name:        "title",
				Value:       "LDAP Proxy",
				Usage:       "Title displayed on the login",
				EnvVars:     []string{"LDAP_PROXY_TITLE"},
				Destination: &config.Server.Title,
			},
			&cli.StringFlag{
				Name:        "endpoint",
				Value:       "",
				Usage:       "Endpoint to proxy requests to",
				EnvVars:     []string{"LDAP_PROXY_ENDPOINT"},
				Destination: &config.Server.Endpoint,
			},
			&cli.StringFlag{
				Name:        "ldap-address",
				Value:       "ldap:389",
				Usage:       "Hostname of the LDAP server",
				EnvVars:     []string{"LDAP_PROXY_SERVER_ADDRESS"},
				Destination: &config.LDAP.Addr,
			},
			&cli.StringFlag{
				Name:        "ldap-bind-username",
				Value:       "",
				Usage:       "Username for bind to server",
				EnvVars:     []string{"LDAP_PROXY_BIND_USERNAME"},
				Destination: &config.LDAP.BindUsername,
			},
			&cli.StringFlag{
				Name:        "ldap-bind-password",
				Value:       "",
				Usage:       "Password for bind to server",
				EnvVars:     []string{"LDAP_PROXY_BIND_PASSWORD"},
				Destination: &config.LDAP.BindPassword,
			},
			&cli.StringFlag{
				Name:        "ldap-base-dn",
				Value:       "",
				Usage:       "Base DN for LDAP server",
				EnvVars:     []string{"LDAP_PROXY_BASE_DN"},
				Destination: &config.LDAP.BaseDN,
			},
			&cli.StringFlag{
				Name:        "ldap-filter-dn",
				Value:       "(&(objectClass=person)(sAMAccountName={login}))",
				Usage:       "User filter for LDAP server",
				EnvVars:     []string{"LDAP_PROXY_FILTER_DN"},
				Destination: &config.LDAP.FilterDN,
			},
			&cli.StringFlag{
				Name:        "ldap-user-attr",
				Value:       "sAMAccountName",
				Usage:       "Attribute for username",
				EnvVars:     []string{"LDAP_PROXY_USER_ATTR"},
				Destination: &config.LDAP.UserAttr,
			},
			&cli.StringFlag{
				Name:        "ldap-user-header",
				Value:       "X-PROXY-USER",
				Usage:       "Header for username",
				EnvVars:     []string{"LDAP_PROXY_USER_HEADER"},
				Destination: &config.LDAP.UserHeader,
			},
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			logrus.Infof("Starting on %s", config.Server.Addr)

			cfg, err := ssl()

			if err != nil {
				return err
			}

			server := &http.Server{
				Addr:         config.Server.Addr,
				Handler:      router.Load(),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				TLSConfig:    cfg,
			}

			if err := startServer(server); err != nil {
				logrus.Fatal(err)
			}

			return nil
		},
	}
}

func curves() []tls.CurveID {
	return []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
}

func ciphers() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
}

func ssl() (*tls.Config, error) {
	if config.Server.Cert != "" && config.Server.Key != "" {
		cert, err := tls.LoadX509KeyPair(
			config.Server.Cert,
			config.Server.Key,
		)

		if err != nil {
			return nil, fmt.Errorf("Failed to load SSL certificates. %s", err)
		}

		return &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         curves(),
			CipherSuites:             ciphers(),
			Certificates:             []tls.Certificate{cert},
		}, nil
	}

	return nil, nil
}
