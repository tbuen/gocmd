// +build release

package config

import (
	"os"
	"path/filepath"
)

var filenameApps string
var filenameBookmarks string
var filenameTabs string

func init() {
	config, err := os.UserConfigDir()
	if err == nil {
		filenameApps = config + string(filepath.Separator) + "gocmd" + string(filepath.Separator) + "apps.xml"
		filenameBookmarks = config + string(filepath.Separator) + "gocmd" + string(filepath.Separator) + "bookmarks.xml"
		filenameTabs = config + string(filepath.Separator) + "gocmd" + string(filepath.Separator) + "tabs.xml"
	}
}
