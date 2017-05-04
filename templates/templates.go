package templates

import (
	"html/template"
	"io/ioutil"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/webhippie/ldap-proxy/config"
)

//go:generate fileb0x ab0x.yaml

// Load initializes the template files.
func Load() *template.Template {
	tpls := template.New("")

	for _, file := range FileNames {
		content, err := ReadFile(file)

		if err != nil {
			logrus.Warnf("Failed to read builtin %s template. %s", file, err)
			continue
		}

		tpls.New(
			path.Base(file),
		).Parse(
			string(content),
		)
	}

	if config.Server.Templates != "" {
		if com.IsDir(config.Server.Templates) {
			files, err := com.GetFileListBySuffix(config.Server.Templates, ".html")

			if err != nil {
				logrus.Warnf("Failed to read custom templates. %s", err)
				return tpls
			}

			for _, file := range files {
				content, err := ioutil.ReadFile(file)

				if err != nil {
					logrus.Warnf("Failed to read custom %s template. %s", file, err)
					continue
				}

				tpls.New(
					path.Base(file),
				).Parse(
					string(content),
				)
			}
		} else {
			logrus.Warnf("Custom templates directory doesn't exist")
		}
	}

	return tpls
}
