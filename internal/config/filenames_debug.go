// +build !release

package config

import (
	"path/filepath"
)

var filenameApps string
var filenameTabs string

func init() {
	filenameApps = "." + string(filepath.Separator) + "config" + string(filepath.Separator) + "apps.xml"
	filenameTabs = "." + string(filepath.Separator) + "cache" + string(filepath.Separator) + "tabs.xml"
}
