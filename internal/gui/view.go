package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/pango"
)

var view struct {
	size, time, user, permissions bool
}

var viewIndex = 0

func drawView(context *cairo.Context, layout *pango.Layout, width, height float64) {
	context.SetSourceRGB(0xF2/255.0, 0x6B/255.0, 0x3A/255.0)
	context.Rectangle(5, 5, width-8, height-9)
	context.Fill()
	context.Rectangle(6, 6, width-12, height-11)
	context.Clip()

	text := "<b>view:</b>   "
	var c string
	if view.size {
		c = "\u2611"
	} else {
		c = "\u2610"
	}
	if viewIndex == 0 {
		text += "<b>" + c + " SIZE</b>   "
	} else {
		text += c + " size   "
	}
	if view.time {
		c = "\u2611"
	} else {
		c = "\u2610"
	}
	if viewIndex == 1 {
		text += "<b>" + c + " MODIFICATION TIME</b>   "
	} else {
		text += c + " modification time   "
	}
	if view.user {
		c = "\u2611"
	} else {
		c = "\u2610"
	}
	if viewIndex == 2 {
		text += "<b>" + c + " USER/GROUP</b>   "
	} else {
		text += c + " user/group   "
	}
	if view.permissions {
		c = "\u2611"
	} else {
		c = "\u2610"
	}
	if viewIndex == 3 {
		text += "<b>" + c + " PERMISSIONS</b>"
	} else {
		text += c + " permissions"
	}

	context.SetSourceRGB(1, 1, 1)
	context.MoveTo(10, 8)
	layout.SetMarkup(text, -1)
	pango.CairoShowLayout(context, layout)
	layout.SetMarkup("", -1)
}

func keyView(key uint) {
	switch key {
	case gdk.KEY_h, gdk.KEY_Left:
		if viewIndex > 0 {
			viewIndex--
		}
		Refresh()
	case gdk.KEY_l, gdk.KEY_Right:
		if viewIndex < 3 {
			viewIndex++
		}
		Refresh()
	case gdk.KEY_j, gdk.KEY_Down, gdk.KEY_k, gdk.KEY_Up, gdk.KEY_space:
		if viewIndex == 0 {
			view.size = !view.size
		} else if viewIndex == 1 {
			view.time = !view.time
		} else if viewIndex == 2 {
			view.user = !view.user
		} else if viewIndex == 3 {
			view.permissions = !view.permissions
		}
		Refresh()
	case gdk.KEY_v, gdk.KEY_q, gdk.KEY_Escape, gdk.KEY_Return:
		mode = MODE_NORMAL
		Refresh()
	}
}
