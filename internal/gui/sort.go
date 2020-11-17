package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend"
	"github.com/tbuen/gocmd/internal/config"
)

func drawSort(context *cairo.Context, layout *pango.Layout, width, height float64, dir *backend.Directory) {
	sortKey, sortOrder := dir.SortKey()

	setSourceColor(context, "F26B3A")
	context.Rectangle(5, 5, width-8, height-9)
	context.Fill()
	context.Rectangle(6, 6, width-12, height-11)
	context.Clip()

	text := "<b>sort:</b>   "
	if sortKey == config.SORT_BY_NAME {
		if sortOrder == config.SORT_ASCENDING {
			text += "<b>NAME\u2191</b>   "
		} else {
			text += "<b>NAME\u2193</b>   "
		}
	} else {
		text += "name    "
	}
	if sortKey == config.SORT_BY_EXT {
		if sortOrder == config.SORT_ASCENDING {
			text += "<b>EXTENSION\u2191</b>   "
		} else {
			text += "<b>EXTENSION\u2193</b>   "
		}
	} else {
		text += "extension    "
	}
	if sortKey == config.SORT_BY_SIZE {
		if sortOrder == config.SORT_ASCENDING {
			text += "<b>SIZE\u2191</b>   "
		} else {
			text += "<b>SIZE\u2193</b>   "
		}
	} else {
		text += "size    "
	}
	if sortKey == config.SORT_BY_TIME {
		if sortOrder == config.SORT_ASCENDING {
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
	dir := backend.GetDirectory(backend.PANEL_ACTIVE)
	sortKey, sortOrder := dir.SortKey()
	switch key {
	case gdk.KEY_h, gdk.KEY_Left:
		if sortKey == config.SORT_BY_EXT {
			sortKey = config.SORT_BY_NAME
		} else if sortKey == config.SORT_BY_SIZE {
			sortKey = config.SORT_BY_EXT
		} else if sortKey == config.SORT_BY_TIME {
			sortKey = config.SORT_BY_SIZE
		}
		dir.SetSortKey(sortKey, sortOrder)
	case gdk.KEY_l, gdk.KEY_Right:
		if sortKey == config.SORT_BY_NAME {
			sortKey = config.SORT_BY_EXT
		} else if sortKey == config.SORT_BY_EXT {
			sortKey = config.SORT_BY_SIZE
		} else if sortKey == config.SORT_BY_SIZE {
			sortKey = config.SORT_BY_TIME
		}
		dir.SetSortKey(sortKey, sortOrder)
	case gdk.KEY_j, gdk.KEY_Down:
		dir.SetSortKey(sortKey, config.SORT_DESCENDING)
	case gdk.KEY_k, gdk.KEY_Up:
		dir.SetSortKey(sortKey, config.SORT_ASCENDING)
	case gdk.KEY_s, gdk.KEY_q, gdk.KEY_Escape, gdk.KEY_Return:
		mode = MODE_NORMAL
		Refresh()
	}
}
