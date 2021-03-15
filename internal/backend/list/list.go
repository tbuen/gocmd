package list

import (
	"github.com/tbuen/gocmd/internal/backend/gui"
)

type List struct {
	list      []interface{}
	selection int
	offset    int
}

func (l *List) Elements() []interface{} {
	return l.list
}

func (l *List) Element() (e interface{}) {
	if l.selection < len(l.list) {
		e = l.list[l.selection]
	}
	return
}

func (l *List) Selection() int {
	return l.selection
}

func (l *List) SetSelectionRel(n int) {
	if n >= 0 {
		l.SetSelection(l.selection + n)
	} else {
		n = -n
		if n > l.selection {
			n = l.selection
		}
		l.SetSelection(l.selection - n)
	}
}

func (l *List) SetSelection(n int) {
	if len(l.list) == 0 {
		l.selection = 0
	} else {
		l.selection = n
		if l.selection < 0 || l.selection >= len(l.list) {
			l.selection = len(l.list) - 1
		}
	}
	gui.Refresh()
}

func (l *List) Offset() int {
	return l.offset
}

func (l *List) SetOffset(offset int) {
	l.offset = offset
}
