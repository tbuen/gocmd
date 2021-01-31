package list

import (
	"github.com/tbuen/gocmd/internal/config"
	. "github.com/tbuen/gocmd/internal/global"
)

type Bookmarks struct {
	bookmarks []Bookmark
	lst       list
}

func (l *Bookmarks) Bookmarks() []Bookmark {
	l.bookmarks = config.Bookmarks().Get()
	l.lst.setSelRel(len(l.bookmarks), 0)
	return l.bookmarks
}

func (l *Bookmarks) Selection() int {
	return l.lst.selection
}

func (l *Bookmarks) SetSelectionRelative(n int) {
	l.lst.setSelRel(len(l.bookmarks), n)
}

func (l *Bookmarks) SetSelectionAbsolute(n int) {
	l.lst.setSelAbs(len(l.bookmarks), n)
}

func (l *Bookmarks) DispOffset() int {
	return l.lst.offset
}

func (l *Bookmarks) SetDispOffset(offset int) {
	l.lst.offset = offset
}
