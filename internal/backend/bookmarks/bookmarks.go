package bookmarks

import (
	"github.com/tbuen/gocmd/internal/backend/list"
	"github.com/tbuen/gocmd/internal/config"
)

type Bookmarks struct {
	list.List
}

func (b *Bookmarks) Get() []config.Bookmark {
	ee := b.Elements()
	bb := make([]config.Bookmark, len(ee))
	for i, e := range ee {
		bb[i] = e.(config.Bookmark)
	}
	return bb
}
