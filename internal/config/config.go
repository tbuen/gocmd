package config

import (
	"github.com/tbuen/gocmd/internal/config/apps"
	"github.com/tbuen/gocmd/internal/config/bookmarks"
	"github.com/tbuen/gocmd/internal/config/tabs"
)

var cfg struct {
	apps      apps.Config
	bookmarks bookmarks.Config
	tabs      tabs.Config
}

func Load() {
	cfg.apps.Load(filenameApps)
	cfg.bookmarks.Load(filenameBookmarks)
	cfg.tabs.Load(filenameTabs)
}

func Save() {
	cfg.bookmarks.Save(filenameBookmarks)
	cfg.tabs.Save(filenameTabs)
}

func Apps() *apps.Config {
	return &cfg.apps
}

func Bookmarks() *bookmarks.Config {
	return &cfg.bookmarks
}
