package assets

import (
	"net/http"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/webhippie/ldap-proxy/pkg/config"

	// Dummy import to have the dep managed.
	_ "golang.org/x/net/webdav"
)

//go:generate fileb0x ab0x.yaml

// Load initializes the static files.
func Load(cfg *config.Config) http.FileSystem {
	return ChainedFS{
		cfg: cfg,
	}
}

// ChainedFS is a simple HTTP filesystem including custom path.
type ChainedFS struct {
	cfg *config.Config
}

// Open just implements the HTTP filesystem interface.
func (c ChainedFS) Open(origPath string) (http.File, error) {
	if c.cfg.Server.Assets != "" {
		if stat, err := os.Stat(c.cfg.Server.Assets); err == nil && stat.IsDir() {
			customPath := path.Join(
				c.cfg.Server.Assets,
				origPath,
			)

			if _, err := os.Stat(customPath); !os.IsNotExist(err) {
				f, err := os.Open(customPath)

				if err != nil {
					return nil, err
				}

				return f, nil
			}
		} else {
			log.Warn().
				Str("dir", c.cfg.Server.Assets).
				Msg("custom assets directory doesn't exist")
		}
	}

	f, err := FS.OpenFile(CTX, origPath, os.O_RDONLY, 0644)

	if err != nil {
		return nil, err
	}

	return f, nil
}
