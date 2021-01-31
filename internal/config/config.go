package config

import (
	"github.com/tbuen/gocmd/internal/config/bookmarks"
)

var bookmarkCfg bookmarks.Config

func Load() {
	bookmarkCfg.Load(filenameBookmarks)
}

func Save() {
	bookmarkCfg.Save(filenameBookmarks)
}

func Bookmarks() *bookmarks.Config {
	return &bookmarkCfg
}
