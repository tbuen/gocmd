package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend/dir"
	"github.com/tbuen/gocmd/internal/backend/panel"
	. "github.com/tbuen/gocmd/internal/global"
)

func drawSort(context *cairo.Context, layout *pango.Layout, width, height float64, d *dir.Directory) {
	sortKey, sortOrder := d.SortKey()

	setSourceColor(context, "F26B3A")
	context.Rectangle(5, 5, width-8, height-9)
	context.Fill()
	context.Rectangle(6, 6, width-12, height-11)
	context.Clip()

	text := "<b>sort:</b>   "
	if sortKey == SORT_BY_NAME {
		if sortOrder == SORT_ASCENDING {
			text += "<b>NAME\u2191</b>   "
		} else {
			text += "<b>NAME\u2193</b>   "
		}
	} else {
		text += "name    "
	}
	if sortKey == SORT_BY_EXT {
		if sortOrder == SORT_ASCENDING {
			text += "<b>EXTENSION\u2191</b>   "
		} else {
			text += "<b>EXTENSION\u2193</b>   "
		}
	} else {
		text += "extension    "
	}
	if sortKey == SORT_BY_SIZE {
		if sortOrder == SORT_ASCENDING {
			text += "<b>SIZE\u2191</b>   "
		} else {
			text += "<b>SIZE\u2193</b>   "
		}
	} else {
		text += "size    "
	}
	if sortKey == SORT_BY_TIME {
		if sortOrder == SORT_ASCENDING {
			text += "<b>MODIFICATION TIME\u2191</b>   "
		} else {
			text += "<b>MODIFICATION TIME\u2193</b>   "
		}
	} else {
		text += "modification time"
	}

	setSourceColor(context, "FFFFFF")
	context.MoveTo(10, 8)
	layout.SetMarkup(text, -1)
	pango.CairoShowLayout(context, layout)
	layout.SetMarkup("", -1)
}

func keySort(key uint) {
	d := panel.Active().Tab().Directory()
	sortKey, sortOrder := d.SortKey()
	switch key {
	case gdk.KEY_h, gdk.KEY_Left:
		if sortKey == SORT_BY_EXT {
			sortKey = SORT_BY_NAME
		} else if sortKey == SORT_BY_SIZE {
			sortKey = SORT_BY_EXT
		} else if sortKey == SORT_BY_TIME {
			sortKey = SORT_BY_SIZE
		}
		d.SetSortKey(sortKey, sortOrder)
	case gdk.KEY_l, gdk.KEY_Right:
		if sortKey == SORT_BY_NAME {
			sortKey = SORT_BY_EXT
		} else if sortKey == SORT_BY_EXT {
			sortKey = SORT_BY_SIZE
		} else if sortKey == SORT_BY_SIZE {
			sortKey = SORT_BY_TIME
		}
		d.SetSortKey(sortKey, sortOrder)
	case gdk.KEY_j, gdk.KEY_Down:
		d.SetSortKey(sortKey, SORT_DESCENDING)
	case gdk.KEY_k, gdk.KEY_Up:
		d.SetSortKey(sortKey, SORT_ASCENDING)
	case gdk.KEY_s, gdk.KEY_q, gdk.KEY_Escape, gdk.KEY_Return:
		mode = MODE_NORMAL
		Refresh()
	}
}
