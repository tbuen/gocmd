package tab

import (
	//"container/list"
	"github.com/tbuen/gocmd/internal/backend/dir"
	"github.com/tbuen/gocmd/internal/backend/gui"
	"github.com/tbuen/gocmd/internal/backend/listing"
	//"github.com/tbuen/gocmd/internal/config"
	//. "github.com/tbuen/gocmd/internal/global"
	//"github.com/tbuen/gocmd/internal/log"
	//"path/filepath"
)

const (
	MODE_DIRECTORY = iota
	MODE_BOOKMARKS
)

type Tab struct {
	mode int
	dir  *dir.Directory
	bl   listing.Bookmarks
}

/*type Tabs struct {
	Panel  int
	Titles []string
	Active int
	Offset float64
}*/

func New() *Tab {
	t := new(Tab)
	t.dir = dir.New()
	return t
}

func (t *Tab) Clone() *Tab {
	c := new(Tab)
	//c.dir = t.dir.Clone()
	// TODO clone := p.ActiveTab().Directory().Clone()
	/*src := p.Tab().Directory()
	sortKey, sortOrder := src.SortKey()
	clone := newDirectory(src.Path(), sortKey, sortOrder, src.Hidden())
	p.insert(clone)*/
	return c
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

func (t *Tab) Directory() (d *dir.Directory) {
	if t.mode == MODE_DIRECTORY {
		d = t.dir
	}
	return
}

func (t *Tab) Bookmarks() (b *listing.Bookmarks) {
	if t.mode == MODE_BOOKMARKS {
		b = &t.bl
	}
	return
}
