// +build !release

package config

import (
	"path/filepath"
)

var filenameApps string
var filenameBookmarks string
var filenameTabs string

func init() {
	filenameApps = "." + string(filepath.Separator) + "config" + string(filepath.Separator) + "apps.xml"
	filenameBookmarks = "." + string(filepath.Separator) + "config" + string(filepath.Separator) + "bookmarks.xml"
	filenameTabs = "." + string(filepath.Separator) + "config" + string(filepath.Separator) + "tabs.xml"
}
