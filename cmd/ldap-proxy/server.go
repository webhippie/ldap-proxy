package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/vulcand/oxy/buffer"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"github.com/webhippie/ldap-proxy/pkg/router"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/urfave/cli.v2"
)

var (
	httpsAddr  = "0.0.0.0:443"
	httpAddr   = "0.0.0.0:80"
	healthAddr = "127.0.0.1:9000"
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "server",
		Usage:  "start the integrated server",
		Flags:  serverFlags(cfg),
		Before: serverBefore(cfg),
		Action: serverAction(cfg),
	}
}

func serverFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "health-addr",
			Value:       healthAddr,
			Usage:       "address for metrics and health",
			EnvVars:     []string{"LDAP_PROXY_HEALTH_ADDR"},
			Destination: &cfg.Server.Health,
		},
		&cli.StringFlag{
			Name:        "secure-addr",
			Value:       httpsAddr,
			Usage:       "https address to bind the server",
			EnvVars:     []string{"LDAP_PROXY_SERVER_HTTPS"},
			Destination: &cfg.Server.Secure,
		},
		&cli.StringFlag{
			Name:        "server-addr",
			Value:       httpAddr,
			Usage:       "http address to bind the server",
			EnvVars:     []string{"LDAP_PROXY_SERVER_ADDR"},
			Destination: &cfg.Server.Public,
		},
		&cli.StringFlag{
			Name:        "server-root",
			Value:       "/ldap-proxy",
			Usage:       "root path of the proxy",
			EnvVars:     []string{"LDAP_PROXY_SERVER_ROOT"},
			Destination: &cfg.Server.Root,
		},
		&cli.StringFlag{
			Name:        "server-host",
			Value:       "http://localhost",
			Usage:       "external access to server",
			EnvVars:     []string{"LDAP_PROXY_SERVER_HOST"},
			Destination: &cfg.Server.Host,
		},
		&cli.StringFlag{
			Name:        "server-cert",
			Value:       "",
			Usage:       "path to ssl cert",
			EnvVars:     []string{"LDAP_PROXY_SERVER_CERT"},
			Destination: &cfg.Server.Cert,
		},
		&cli.StringFlag{
			Name:        "server-key",
			Value:       "",
			Usage:       "path to ssl key",
			EnvVars:     []string{"LDAP_PROXY_SERVER_KEY"},
			Destination: &cfg.Server.Key,
		},
		&cli.BoolFlag{
			Name:        "server-autocert",
			Value:       false,
			Usage:       "enable let's encrypt",
			EnvVars:     []string{"LDAP_PROXY_AUTO_CERT"},
			Destination: &cfg.Server.AutoCert,
		},
		&cli.BoolFlag{
			Name:        "strict-curves",
			Value:       false,
			Usage:       "use strict ssl curves",
			EnvVars:     []string{"LDAP_PROXY_STRICT_CURVES"},
			Destination: &cfg.Server.StrictCurves,
		},
		&cli.BoolFlag{
			Name:        "strict-ciphers",
			Value:       false,
			Usage:       "use strict ssl ciphers",
			EnvVars:     []string{"LDAP_PROXY_STRICT_CIPHERS"},
			Destination: &cfg.Server.StrictCiphers,
		},
		&cli.StringFlag{
			Name:        "templates-path",
			Value:       "",
			Usage:       "path to custom templates",
			EnvVars:     []string{"LDAP_PROXY_SERVER_TEMPLATES"},
			Destination: &cfg.Server.Templates,
		},
		&cli.StringFlag{
			Name:        "assets-path",
			Value:       "",
			Usage:       "path to custom assets",
			EnvVars:     []string{"LDAP_PROXY_SERVER_ASSETS"},
			Destination: &cfg.Server.Assets,
		},
		&cli.StringFlag{
			Name:        "storage-path",
			Value:       "storage/",
			Usage:       "folder for storing certs and misc files",
			EnvVars:     []string{"LDAP_PROXY_SERVER_STORAGE"},
			Destination: &cfg.Server.Storage,
		},
		&cli.StringFlag{
			Name:        "proxy-title",
			Value:       "LDAP Proxy",
			Usage:       "title displayed on the login",
			EnvVars:     []string{"LDAP_PROXY_SERVER_TITLE"},
			Destination: &cfg.Proxy.Title,
		},
		&cli.StringSliceFlag{
			Name:    "proxy-endpoint",
			Value:   cli.NewStringSlice(),
			Usage:   "endpoints to proxy requests to",
			EnvVars: []string{"LDAP_PROXY_SERVER_ENDPOINTS"},
		},
		&cli.StringFlag{
			Name:        "user-header",
			Value:       "X-PROXY-USER",
			Usage:       "header for username",
			EnvVars:     []string{"LDAP_PROXY_USER_HEADER"},
			Destination: &cfg.Proxy.UserHeader,
		},
		&cli.StringFlag{
			Name:        "ldap-address",
			Value:       "ldap:389",
			Usage:       "hostname of the ldap server",
			EnvVars:     []string{"LDAP_PROXY_SERVER_ADDRESS"},
			Destination: &cfg.LDAP.Addr,
		},
		&cli.StringFlag{
			Name:        "ldap-username",
			Value:       "",
			Usage:       "username for bind to server",
			EnvVars:     []string{"LDAP_PROXY_USERNAME"},
			Destination: &cfg.LDAP.BindUsername,
		},
		&cli.StringFlag{
			Name:        "ldap-password",
			Value:       "",
			Usage:       "password for bind to server",
			EnvVars:     []string{"LDAP_PROXY_PASSWORD"},
			Destination: &cfg.LDAP.BindPassword,
		},
		&cli.StringFlag{
			Name:        "ldap-base",
			Value:       "",
			Usage:       "base dn for ldap server",
			EnvVars:     []string{"LDAP_PROXY_BASE"},
			Destination: &cfg.LDAP.BaseDN,
		},
		&cli.StringFlag{
			Name:        "ldap-filter",
			Value:       "(&(objectClass=person)(sAMAccountName={login}))",
			Usage:       "user filter for ldap server",
			EnvVars:     []string{"LDAP_PROXY_FILTER"},
			Destination: &cfg.LDAP.FilterDN,
		},
		&cli.StringFlag{
			Name:        "ldap-userattr",
			Value:       "sAMAccountName",
			Usage:       "attribute for username",
			EnvVars:     []string{"LDAP_PROXY_USER_ATTR"},
			Destination: &cfg.LDAP.UserAttr,
		},
	}
}

func serverBefore(cfg *config.Config) cli.BeforeFunc {
	return func(c *cli.Context) error {
		if len(c.StringSlice("proxy-endpoint")) > 0 {
			// StringSliceFlag doesn't support Destination
			cfg.Proxy.Endpoints = c.StringSlice("proxy-endpoint")
		}

		return nil
	}
}

func serverAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		fwd, err := forward.New(
			forward.PassHostHeader(true),
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to initialize forwarder")

			return err
		}

		lb, err := roundrobin.New(fwd)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to initialize balancer")

			return err
		}

		proxy, err := buffer.New(
			lb,
			buffer.Retry(`IsNetworkError() && Attempts() < 3`),
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to initialize buffer")

			return err
		}

		for _, endpoint := range cfg.Proxy.Endpoints {
			parsed, err := url.Parse(endpoint)

			if err != nil {
				log.Warn().
					Err(err).
					Str("endpoint", endpoint).
					Msg("failed to parse endpoint")

				continue
			}

			lb.UpsertServer(parsed)
		}

		var gr run.Group

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

		{
			server := &http.Server{
				Addr:         cfg.Server.Health,
				Handler:      router.Status(cfg),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Health).
					Msg("starting status server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Info().
						Err(err).
						Msg("failed to stop status server gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("status server stopped gracefully")
			})
		}

		if cfg.Server.AutoCert {
			parsed, err := url.Parse(
				cfg.Server.Host,
			)

			if err != nil {
				log.Info().
					Err(err).
					Msg("failed to parse host")

				return err
			}

			manager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(parsed.Host),
				Cache:      autocert.DirCache(path.Join(cfg.Server.Storage, "certs")),
			}

			{
				server := &http.Server{
					Addr:         httpAddr,
					Handler:      router.Redirect(cfg),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", httpAddr).
						Msg("starting http server")

					return server.ListenAndServe()
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop http server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("http server stopped gracefully")
				})
			}

			{
				server := &http.Server{
					Addr:         httpsAddr,
					Handler:      router.Load(cfg, proxy),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					TLSConfig: &tls.Config{
						PreferServerCipherSuites: true,
						MinVersion:               tls.VersionTLS12,
						CurvePreferences:         curves(cfg),
						CipherSuites:             ciphers(cfg),
						GetCertificate:           manager.GetCertificate,
					},
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", httpsAddr).
						Msg("starting https server")

					return server.ListenAndServeTLS("", "")
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop https server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("https server stopped gracefully")
				})
			}

			return gr.Run()
		} else if cfg.Server.Cert != "" && cfg.Server.Key != "" {
			cert, err := tls.LoadX509KeyPair(
				cfg.Server.Cert,
				cfg.Server.Key,
			)

			if err != nil {
				log.Info().
					Err(err).
					Msg("failed to load certificates")

				return err
			}

			{
				server := &http.Server{
					Addr:         cfg.Server.Public,
					Handler:      router.Redirect(cfg),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", cfg.Server.Public).
						Msg("starting http server")

					return server.ListenAndServe()
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop http server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("http server stopped gracefully")
				})
			}

			{
				server := &http.Server{
					Addr:         cfg.Server.Secure,
					Handler:      router.Load(cfg, proxy),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					TLSConfig: &tls.Config{
						PreferServerCipherSuites: true,
						MinVersion:               tls.VersionTLS12,
						CurvePreferences:         curves(cfg),
						CipherSuites:             ciphers(cfg),
						Certificates:             []tls.Certificate{cert},
					},
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", cfg.Server.Secure).
						Msg("starting https server")

					return server.ListenAndServeTLS("", "")
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop https server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("https server stopped gracefully")
				})
			}

			return gr.Run()
		}

		{
			server := &http.Server{
				Addr:         cfg.Server.Public,
				Handler:      router.Load(cfg, proxy),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Public).
					Msg("starting http server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Info().
						Err(err).
						Msg("failed to stop http server gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("http server stopped gracefully")
			})
		}

		return gr.Run()
	}
}

func curves(cfg *config.Config) []tls.CurveID {
	if cfg.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

func ciphers(cfg *config.Config) []uint16 {
	if cfg.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}
