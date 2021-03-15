package tab

import (
	"github.com/tbuen/gocmd/internal/backend/bookmarks"
	"github.com/tbuen/gocmd/internal/backend/dir"
	"github.com/tbuen/gocmd/internal/backend/gui"
	"github.com/tbuen/gocmd/internal/config"
)

const (
	MODE_DIRECTORY = iota
	MODE_BOOKMARKS
)

type Tab struct {
	mode  int
	dir   *dir.Directory
	bookm bookmarks.Bookmarks
}

func New() *Tab {
	t := new(Tab)
	t.dir = dir.New()
	return t
}

func NewWithConfig(cfg config.Directory) *Tab {
	t := new(Tab)
	t.dir = dir.NewWithConfig(cfg)
	return t
}

func (t *Tab) Clone() *Tab {
	c := new(Tab)
	c.dir = t.dir.Clone()
	return c
}

func (t *Tab) Destroy() {
	t.dir.Destroy()
}

func (t *Tab) Config() config.Directory {
	return t.dir.Config()
}

func (t *Tab) Reload() {
	t.dir.Reload()
}

func (t *Tab) Mode() int {
	return t.mode
}

func (t *Tab) ShowBookmarks() {
	if t.mode == MODE_DIRECTORY {
		t.mode = MODE_BOOKMARKS
		gui.Refresh()
	}
}

func (t *Tab) HideBookmarks() {
	if t.mode == MODE_BOOKMARKS {
		t.mode = MODE_DIRECTORY
		gui.Refresh()
	}
}

func (t *Tab) Directory() *dir.Directory {
	return t.dir
}

func (t *Tab) Bookmarks() *bookmarks.Bookmarks {
	return &t.bookm
}
