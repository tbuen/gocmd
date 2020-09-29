// +build !release

package config

import (
	"path/filepath"
)

var filenameApps string

func init() {
	filenameApps = "." + string(filepath.Separator) + "configs" + string(filepath.Separator) + "apps.xml"
}
