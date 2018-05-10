package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/webhippie/ldap-proxy/pkg/config"
	"gopkg.in/urfave/cli.v2"
)

// Health provides the sub-command to perform a health check.
func Health(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "health",
		Usage:  "perform health checks for service",
		Flags:  healthFlags(cfg),
		Action: healthAction(cfg),
	}
}

func healthFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "health-addr",
			Value:       healthAddr,
			Usage:       "health address for service",
			EnvVars:     []string{"LDAP_PROXY_HEALTH_ADDR"},
			Destination: &cfg.Server.Health,
		},
	}
}

func healthAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		resp, err := http.Get(
			fmt.Sprintf(
				"http://%s/healthz",
				cfg.Server.Health,
			),
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to request health check")

			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Error().
				Err(err).
				Msg("health seems to be in a bad state")

			return err
		}

		return nil
	}
}
