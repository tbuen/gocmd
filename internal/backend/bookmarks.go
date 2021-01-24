package backend

type Bookmark struct {
	path string
}

type Bookmarks struct {
	bookmarks  []Bookmark
	selection  int
	dispOffset int
}

func newBookmarks() (bookmarks *Bookmarks) {
	bookmarks = new(Bookmarks)
	bookmarks.bookmarks = append(bookmarks.bookmarks, Bookmark{"aaa"})
	bookmarks.bookmarks = append(bookmarks.bookmarks, Bookmark{"bbb"})
	bookmarks.bookmarks = append(bookmarks.bookmarks, Bookmark{"ccc"})
	return
}

func (bm *Bookmark) Path() string {
	return bm.path
}

func (bm *Bookmarks) Bookmarks() []Bookmark {
	return bm.bookmarks
}

func (bm *Bookmarks) Selection() int {
	return bm.selection
}

func (bm *Bookmarks) SetSelectionRelative(n int) {
	if n > 0 {
		bm.SetSelectionAbsolute(bm.selection + n)
	} else {
		n = -n
		if n > bm.selection {
			n = bm.selection
		}
		bm.SetSelectionAbsolute(bm.selection - n)
	}
}

func (bm *Bookmarks) SetSelectionAbsolute(n int) {
	bm.selection = n
	if bm.selection < 0 || bm.selection >= len(bm.bookmarks) {
		bm.selection = len(bm.bookmarks) - 1
	}
	guiRefresh()
}

func (bm *Bookmarks) DispOffset() int {
	return bm.dispOffset
}

func (bm *Bookmarks) SetDispOffset(offset int) {
	bm.dispOffset = offset
}
