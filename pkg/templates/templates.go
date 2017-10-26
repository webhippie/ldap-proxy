package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/webhippie/ldap-proxy/pkg/config"
)

//go:generate retool -tool-dir ../../_tools do fileb0x ab0x.yaml

// File returns the content of a specifc template.
func File(name string) ([]byte, error) {
	name = strings.TrimPrefix(
		name,
		"templates/",
	)

	if config.Server.Templates != "" {
		if com.IsDir(config.Server.Templates) {
			pathName := path.Join(
				config.Server.Templates,
				name,
			)

			if com.IsFile(pathName) {
				return ioutil.ReadFile(pathName)
			}
		} else {
			logrus.Warnf("Custom templates directory doesn't exist")
		}
	}

	return ReadFile(name)
}

// Names returns a list of all available templates.
func Names() []string {
	result := []string{}

	if config.Server.Templates != "" {
		if com.IsDir(config.Server.Templates) {
			files, err := com.GetFileListBySuffix(config.Server.Templates, ".tmpl")

			if err != nil {
				logrus.Warnf("Failed to read custom templates. %s", err)
			} else {
				for _, file := range files {
					result = append(
						result,
						fmt.Sprintf(
							"templates/%s",
							strings.TrimPrefix(
								file,
								config.Server.Templates,
							),
						),
					)
				}
			}
		} else {
			logrus.Warnf("Custom templates directory doesn't exist")
		}
	}

	for _, file := range FileNames {
		result = append(
			result,
			fmt.Sprintf(
				"templates/%s",
				strings.TrimPrefix(file, "./"),
			),
		)
	}

	return result
}

// Funcs provides some general usefule template helpers.
func Funcs() []template.FuncMap {
	return []template.FuncMap{
		{
			"split":    strings.Split,
			"join":     strings.Join,
			"toUpper":  strings.ToUpper,
			"toLower":  strings.ToLower,
			"contains": strings.Contains,
			"replace":  strings.Replace,
		},
	}
}
