package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/router"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/urfave/cli.v2"
)

var (
	defaultAddr = "0.0.0.0:8080"
)

// Server provides the sub-command to start the server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the integrated server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server-host",
				Value:       "http://localhost:8080",
				Usage:       "external access to server",
				EnvVars:     []string{"LDAP_PROXY_SERVER_HOST"},
				Destination: &config.Server.Host,
			},
			&cli.StringFlag{
				Name:        "server-addr",
				Value:       defaultAddr,
				Usage:       "address to bind the server",
				EnvVars:     []string{"LDAP_PROXY_SERVER_ADDR"},
				Destination: &config.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "server-root",
				Value:       "/ldap-proxy",
				Usage:       "root path of the proxy",
				EnvVars:     []string{"LDAP_PROXY_SERVER_ROOT"},
				Destination: &config.Server.Root,
			},
			&cli.BoolFlag{
				Name:        "enable-pprof",
				Value:       false,
				Usage:       "enable pprof debugging server",
				EnvVars:     []string{"LDAP_PROXY_SERVER_PPROF"},
				Destination: &config.Server.Pprof,
			},
			&cli.BoolFlag{
				Name:        "enable-prometheus",
				Value:       false,
				Usage:       "enable prometheus exporter",
				EnvVars:     []string{"LDAP_PROXY_SERVER_PROMETHEUS"},
				Destination: &config.Server.Prometheus,
			},
			&cli.StringFlag{
				Name:        "server-cert",
				Value:       "",
				Usage:       "path to ssl cert",
				EnvVars:     []string{"LDAP_PROXY_SERVER_CERT"},
				Destination: &config.Server.Cert,
			},
			&cli.StringFlag{
				Name:        "server-key",
				Value:       "",
				Usage:       "path to ssl key",
				EnvVars:     []string{"LDAP_PROXY_SERVER_KEY"},
				Destination: &config.Server.Key,
			},
			&cli.BoolFlag{
				Name:        "enable-letsencrypt",
				Value:       false,
				Usage:       "enable let's encrypt ssl",
				EnvVars:     []string{"LDAP_PROXY_SERVER_LETSENCRYPT"},
				Destination: &config.Server.LetsEncrypt,
			},
			&cli.BoolFlag{
				Name:        "strict-curves",
				Value:       false,
				Usage:       "use strict ssl curves",
				EnvVars:     []string{"LDAP_PROXY_STRICT_CURVES"},
				Destination: &config.Server.StrictCurves,
			},
			&cli.BoolFlag{
				Name:        "strict-ciphers",
				Value:       false,
				Usage:       "use strict ssl ciphers",
				EnvVars:     []string{"LDAP_PROXY_STRICT_CIPHERS"},
				Destination: &config.Server.StrictCiphers,
			},
			&cli.StringFlag{
				Name:        "storage-path",
				Value:       "storage/",
				Usage:       "folder for storing certs and misc files",
				EnvVars:     []string{"LDAP_PROXY_SERVER_STORAGE"},
				Destination: &config.Server.Storage,
			},
			&cli.StringFlag{
				Name:        "templates-path",
				Value:       "",
				Usage:       "path to custom templates",
				EnvVars:     []string{"LDAP_PROXY_SERVER_TEMPLATES"},
				Destination: &config.Server.Templates,
			},
			&cli.StringFlag{
				Name:        "assets-path",
				Value:       "",
				Usage:       "path to custom assets",
				EnvVars:     []string{"LDAP_PROXY_SERVER_ASSETS"},
				Destination: &config.Server.Assets,
			},
			&cli.StringFlag{
				Name:        "proxy-title",
				Value:       "LDAP Proxy",
				Usage:       "title displayed on the login",
				EnvVars:     []string{"LDAP_PROXY_SERVER_TITLE"},
				Destination: &config.Server.Title,
			},
			&cli.StringFlag{
				Name:        "proxy-endpoint",
				Value:       "",
				Usage:       "endpoint to proxy requests to",
				EnvVars:     []string{"LDAP_PROXY_SERVER_ENDPOINT"},
				Destination: &config.Server.Endpoint,
			},
			&cli.StringFlag{
				Name:        "user-header",
				Value:       "X-PROXY-USER",
				Usage:       "header for username",
				EnvVars:     []string{"LDAP_PROXY_USER_HEADER"},
				Destination: &config.LDAP.UserHeader,
			},
			&cli.StringFlag{
				Name:        "ldap-address",
				Value:       "ldap:389",
				Usage:       "hostname of the ldap server",
				EnvVars:     []string{"LDAP_PROXY_SERVER_ADDRESS"},
				Destination: &config.LDAP.Addr,
			},
			&cli.StringFlag{
				Name:        "ldap-bind-username",
				Value:       "",
				Usage:       "username for bind to server",
				EnvVars:     []string{"LDAP_PROXY_BIND_USERNAME"},
				Destination: &config.LDAP.BindUsername,
			},
			&cli.StringFlag{
				Name:        "ldap-bind-password",
				Value:       "",
				Usage:       "password for bind to server",
				EnvVars:     []string{"LDAP_PROXY_BIND_PASSWORD"},
				Destination: &config.LDAP.BindPassword,
			},
			&cli.StringFlag{
				Name:        "ldap-base-dn",
				Value:       "",
				Usage:       "base dn for ldap server",
				EnvVars:     []string{"LDAP_PROXY_BASE_DN"},
				Destination: &config.LDAP.BaseDN,
			},
			&cli.StringFlag{
				Name:        "ldap-filter-dn",
				Value:       "(&(objectClass=person)(sAMAccountName={login}))",
				Usage:       "user filter for ldap server",
				EnvVars:     []string{"LDAP_PROXY_FILTER_DN"},
				Destination: &config.LDAP.FilterDN,
			},
			&cli.StringFlag{
				Name:        "ldap-user-attr",
				Value:       "sAMAccountName",
				Usage:       "attribute for username",
				EnvVars:     []string{"LDAP_PROXY_USER_ATTR"},
				Destination: &config.LDAP.UserAttr,
			},
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

			switch strings.ToLower(config.LogLevel) {
			case "debug":
				logger = level.NewFilter(logger, level.AllowDebug())
			case "warn":
				logger = level.NewFilter(logger, level.AllowWarn())
			case "error":
				logger = level.NewFilter(logger, level.AllowError())
			default:
				logger = level.NewFilter(logger, level.AllowInfo())
			}

			logger = log.WithPrefix(logger,
				"app", c.App.Name,
				"ts", log.DefaultTimestampUTC,
			)

			var (
				gr group.Group
			)

			if config.Server.LetsEncrypt || (config.Server.Cert != "" && config.Server.Key != "") {
				cfg, err := ssl(logger)

				if err != nil {
					return err
				}

				if config.Server.LetsEncrypt {
					{
						server := &http.Server{
							Addr:         net.JoinHostPort(addr(), "80"),
							Handler:      redirect(logger),
							ReadTimeout:  5 * time.Second,
							WriteTimeout: 10 * time.Second,
						}

						gr.Add(func() error {
							level.Info(logger).Log(
								"msg", "starting http server",
								"addr", net.JoinHostPort(addr(), "80"),
							)

							return server.ListenAndServe()
						}, func(reason error) {
							ctx, cancel := context.WithTimeout(context.Background(), time.Second)
							defer cancel()

							if err := server.Shutdown(ctx); err != nil {
								level.Error(logger).Log(
									"msg", "failed to shutdown http server gracefully",
									"err", err,
								)

								return
							}

							level.Info(logger).Log(
								"msg", "http server shutdown gracefully",
								"reason", reason,
							)
						})
					}

					{
						server := &http.Server{
							Addr:         net.JoinHostPort(addr(), "443"),
							Handler:      router.Load(logger),
							ReadTimeout:  5 * time.Second,
							WriteTimeout: 10 * time.Second,
							TLSConfig:    cfg,
						}

						gr.Add(func() error {
							level.Info(logger).Log(
								"msg", "starting https server",
								"addr", net.JoinHostPort(addr(), "443"),
							)

							return server.ListenAndServeTLS("", "")
						}, func(reason error) {
							ctx, cancel := context.WithTimeout(context.Background(), time.Second)
							defer cancel()

							if err := server.Shutdown(ctx); err != nil {
								level.Error(logger).Log(
									"msg", "failed to shutdown https server gracefully",
									"err", err,
								)

								return
							}

							level.Info(logger).Log(
								"msg", "https server shutdown gracefully",
								"reason", reason,
							)
						})
					}
				} else {
					{
						server := &http.Server{
							Addr:         config.Server.Addr,
							Handler:      router.Load(logger),
							ReadTimeout:  5 * time.Second,
							WriteTimeout: 10 * time.Second,
							TLSConfig:    cfg,
						}

						gr.Add(func() error {
							level.Info(logger).Log(
								"msg", "starting https server",
								"addr", config.Server.Addr,
							)

							return server.ListenAndServeTLS("", "")
						}, func(reason error) {
							ctx, cancel := context.WithTimeout(context.Background(), time.Second)
							defer cancel()

							if err := server.Shutdown(ctx); err != nil {
								level.Error(logger).Log(
									"msg", "failed to shutdown https server gracefully",
									"err", err,
								)

								return
							}

							level.Info(logger).Log(
								"msg", "https server shutdown gracefully",
								"reason", reason,
							)
						})
					}
				}
			} else {
				{
					server := &http.Server{
						Addr:         config.Server.Addr,
						Handler:      router.Load(logger),
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 10 * time.Second,
					}

					gr.Add(func() error {
						level.Info(logger).Log(
							"msg", "starting http server",
							"addr", config.Server.Addr,
						)

						return server.ListenAndServe()
					}, func(reason error) {
						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()

						if err := server.Shutdown(ctx); err != nil {
							level.Error(logger).Log(
								"msg", "failed to shutdown http server gracefully",
								"err", err,
							)

							return
						}

						level.Info(logger).Log(
							"msg", "http server shutdown gracefully",
							"reason", reason,
						)
					})
				}
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
				})
			}

			return gr.Run()
		},
	}
}

func addr() string {
	splitAddr := strings.SplitN(
		config.Server.Addr,
		":",
		2,
	)

	return splitAddr[0]
}

func curves() []tls.CurveID {
	if config.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

func ciphers() []uint16 {
	if config.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}

func ssl(logger log.Logger) (*tls.Config, error) {
	if config.Server.LetsEncrypt {
		if config.Server.Addr != defaultAddr {
			level.Info(logger).Log(
				"msg", "enabled let's encrypt, overwriting the port",
			)
		}

		parsed, err := url.Parse(
			config.Server.Host,
		)

		if err != nil {
			level.Error(logger).Log(
				"msg", "failed to parse host",
				"err", err,
			)

			return nil, err
		}

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(parsed.Host),
			Cache:      autocert.DirCache(path.Join(config.Server.Storage, "certs")),
		}

		return &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         curves(),
			CipherSuites:             ciphers(),
			GetCertificate:           certManager.GetCertificate,
		}, nil
	}

	if config.Server.Cert != "" && config.Server.Key != "" {
		cert, err := tls.LoadX509KeyPair(
			config.Server.Cert,
			config.Server.Key,
		)

		if err != nil {
			level.Error(logger).Log(
				"msg", "failed to load certificates",
				"err", err,
			)

			return nil, err
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

func redirect(logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := strings.Join(
			[]string{
				"https://",
				r.Host,
				r.URL.Path,
			},
			"",
		)

		if len(r.URL.RawQuery) > 0 {
			target += "?" + r.URL.RawQuery
		}

		level.Debug(logger).Log(
			"msg", "redirecting to https",
			"target", target,
		)

		http.Redirect(w, r, target, http.StatusPermanentRedirect)
	})
}
