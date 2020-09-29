// +build release

package config

import (
	"os"
	"path/filepath"
)

var filenameApps string

func init() {
	config, err := os.UserConfigDir()
	if err == nil {
		filenameApps = config + string(filepath.Separator) + "gocmd" + string(filepath.Separator) + "apps.xml"
	}
}
