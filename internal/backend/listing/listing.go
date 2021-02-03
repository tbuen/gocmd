package listing

import (
	"github.com/tbuen/gocmd/internal/backend/gui"
)

type listing struct {
	selection int
	offset    int
}

func (lst *listing) setSelRel(l, n int) {
	if n >= 0 {
		lst.setSelAbs(l, lst.selection+n)
	} else {
		n = -n
		if n > lst.selection {
			n = lst.selection
		}
		lst.setSelAbs(l, lst.selection-n)
	}
}

func (lst *listing) setSelAbs(l, n int) {
	lst.selection = n
	if lst.selection < 0 || lst.selection >= l {
		lst.selection = l - 1
	}
	gui.Refresh()
}
