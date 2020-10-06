package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend"
)

func drawSort(context *cairo.Context, layout *pango.Layout, width, height float64, dir backend.Directory) {
	crit, desc := dir.Sort()

	context.SetSourceRGB(0xF2/255.0, 0x6B/255.0, 0x3A/255.0)
	context.Rectangle(5, 5, width-8, height-9)
	context.Fill()
	context.Rectangle(6, 6, width-12, height-11)
	context.Clip()

	text := "<b>sort:</b>   "
	if crit == backend.SORT_NAME {
		if desc {
			text += "<b>NAME\u2193</b>   "
		} else {
			text += "<b>NAME\u2191</b>   "
		}
	} else {
		text += "name    "
	}
	if crit == backend.SORT_EXT {
		if desc {
			text += "<b>EXTENSION\u2193</b>   "
		} else {
			text += "<b>EXTENSION\u2191</b>   "
		}
	} else {
		text += "extension    "
	}
	if crit == backend.SORT_SIZE {
		if desc {
			text += "<b>SIZE\u2193</b>   "
		} else {
			text += "<b>SIZE\u2191</b>   "
		}
	} else {
		text += "size    "
	}
	if crit == backend.SORT_TIME {
		if desc {
			text += "<b>MODIFICATION TIME\u2193</b>   "
		} else {
			text += "<b>MODIFICATION TIME\u2191</b>   "
		}
	} else {
		text += "modification time"
	}

	context.SetSourceRGB(1, 1, 1)
	context.MoveTo(10, 8)
	layout.SetMarkup(text, -1)
	pango.CairoShowLayout(context, layout)
	layout.SetMarkup("", -1)
}

func keySort(key uint) {
	dir := backend.GetDirectory(backend.PANEL_ACTIVE)
	crit, desc := dir.Sort()
	switch key {
	case gdk.KEY_h, gdk.KEY_Left:
		if crit == backend.SORT_EXT {
			crit = backend.SORT_NAME
		} else if crit == backend.SORT_SIZE {
			crit = backend.SORT_EXT
		} else if crit == backend.SORT_TIME {
			crit = backend.SORT_SIZE
		}
		dir.SetSort(crit, desc)
	case gdk.KEY_l, gdk.KEY_Right:
		if crit == backend.SORT_NAME {
			crit = backend.SORT_EXT
		} else if crit == backend.SORT_EXT {
			crit = backend.SORT_SIZE
		} else if crit == backend.SORT_SIZE {
			crit = backend.SORT_TIME
		}
		dir.SetSort(crit, desc)
	case gdk.KEY_j, gdk.KEY_Down:
		if !desc {
			desc = true
		}
		dir.SetSort(crit, desc)
	case gdk.KEY_k, gdk.KEY_Up:
		if desc {
			desc = false
		}
		dir.SetSort(crit, desc)
	case gdk.KEY_s, gdk.KEY_q, gdk.KEY_Escape, gdk.KEY_Return:
		mode = MODE_NORMAL
		Refresh()
	}
}
