package templates

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/ldap-proxy/pkg/config"

	// Dummy import to have the dep managed.
	_ "golang.org/x/net/webdav"
)

//go:generate fileb0x ab0x.yaml

// Load initializes the template files.
func Load(cfg *config.Config) *template.Template {
	tpls := template.New(
		"",
	).Funcs(
		sprig.FuncMap(),
	)

	files, err := WalkDirs(
		"",
		false,
	)

	if err != nil {
		log.Warn().
			Err(err).
			Msg("failed to get builtin template list")
	} else {
		for _, name := range files {
			file, readErr := ReadFile(name)

			if readErr != nil {
				log.Warn().
					Err(readErr).
					Str("file", name).
					Msg("failed to read builtin template")
			}

			_, parseErr := tpls.New(
				name,
			).Parse(
				string(file),
			)

			if parseErr != nil {
				log.Warn().
					Err(parseErr).
					Str("file", name).
					Msg("failed to parse builtin template")
			}
		}
	}

	if cfg.Server.Templates != "" {
		if stat, err := os.Stat(cfg.Server.Templates); err == nil && stat.IsDir() {
			files := []string{}

			filepath.Walk(cfg.Server.Templates, func(path string, f os.FileInfo, err error) error {
				if f.IsDir() {
					return nil
				}

				files = append(
					files,
					path,
				)

				return nil
			})

			for _, name := range files {
				file, readErr := ioutil.ReadFile(name)

				if readErr != nil {
					log.Warn().
						Err(readErr).
						Str("file", name).
						Msg("failed to read custom template")
				}

				_, parseErr := tpls.New(
					strings.TrimPrefix(
						strings.TrimPrefix(
							name,
							cfg.Server.Templates,
						),
						"/",
					),
				).Parse(
					string(file),
				)

				if parseErr != nil {
					log.Warn().
						Err(parseErr).
						Str("file", name).
						Msg("failed to parse custom template")
				}
			}
		} else {
			log.Warn().
				Str("dir", cfg.Server.Templates).
				Msg("custom templates directory doesn't exist")
		}
	}

	return tpls
}
