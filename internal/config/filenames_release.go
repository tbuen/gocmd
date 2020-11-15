// +build release

package config

import (
	"os"
	"path/filepath"
)

var filenameApps string
var filenameTabs string

func init() {
	config, err := os.UserConfigDir()
	if err == nil {
		filenameApps = config + string(filepath.Separator) + "gocmd" + string(filepath.Separator) + "apps.xml"
	}
	cache, err := os.UserCacheDir()
	if err == nil {
		filenameTabs = cache + string(filepath.Separator) + "gocmd" + string(filepath.Separator) + "tabs.xml"
	}
}
