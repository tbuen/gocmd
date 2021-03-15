package gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tbuen/gocmd/internal/backend/panel"
	"github.com/tbuen/gocmd/internal/backend/tab"
	"github.com/tbuen/gocmd/internal/log"
)

const (
	MODE_NORMAL = iota
	MODE_SORT
	MODE_VIEW
)

var mode int

func onKeyPress(win *gtk.ApplicationWindow, ev *gdk.Event) {
	keyEvent := gdk.EventKeyNewFromEvent(ev)
	log.Println(log.GUI, "key pressed:", keyEvent.KeyVal())

	switch mode {
	case MODE_NORMAL:
		keyNormal(win, keyEvent.KeyVal())
	case MODE_SORT:
		keySort(keyEvent.KeyVal())
	case MODE_VIEW:
		keyView(keyEvent.KeyVal())
	}
}

func keyNormal(win *gtk.ApplicationWindow, key uint) {
	switch key {
	case gdk.KEY_Q:
		// TODO: Ctrl-Q, Alt-Q etc. should not work...
		win.Close()
	case gdk.KEY_Tab:
		panel.Toggle()
	case gdk.KEY_t:
		panel.Active().New()
	case gdk.KEY_T:
		panel.Active().Clone()
	case gdk.KEY_w:
		panel.Active().Delete()
	case gdk.KEY_h, gdk.KEY_Left:
		panel.Active().Prev()
	case gdk.KEY_H: // TODO Shift+Left
		panel.Active().First()
	case gdk.KEY_l, gdk.KEY_Right:
		panel.Active().Next()
	case gdk.KEY_L: // TODO Shift+Right
		panel.Active().Last()
	default:
		switch panel.Active().Tab().Mode() {
		case tab.MODE_DIRECTORY:
			keyDirectory(key)
		case tab.MODE_BOOKMARKS:
			keyBookmark(key)
		}
	}
}

func keyDirectory(key uint) {
	dir := panel.Active().Tab().Directory()
	switch key {
	case gdk.KEY_r:
		dir.Reload()
	case gdk.KEY_period:
		dir.ToggleHidden()
	case gdk.KEY_j, gdk.KEY_Down:
		dir.SetSelectionRel(1)
	case gdk.KEY_J, gdk.KEY_Page_Down:
		dir.SetSelectionRel(20)
	case gdk.KEY_k, gdk.KEY_Up:
		dir.SetSelectionRel(-1)
	case gdk.KEY_K, gdk.KEY_Page_Up:
		dir.SetSelectionRel(-20)
	case gdk.KEY_g, gdk.KEY_Home:
		dir.SetSelection(0)
	case gdk.KEY_G, gdk.KEY_End:
		dir.SetSelection(-1)
	case gdk.KEY_m:
		dir.ToggleMarkSelected()
	case gdk.KEY_M:
		dir.ToggleMarkAll()
	case gdk.KEY_u, gdk.KEY_numbersign:
		dir.GoUp()
	case gdk.KEY_Return:
		dir.Enter()
	case gdk.KEY_asciicircum:
		dir.Root()
	case gdk.KEY_asciitilde:
		dir.Home()
	case gdk.KEY_F3:
		dir.View()
	case gdk.KEY_F4:
		dir.Edit()
	case gdk.KEY_v:
		mode = MODE_VIEW
		Refresh()
	case gdk.KEY_s:
		mode = MODE_SORT
		Refresh()
	case gdk.KEY_b:
		panel.Active().Tab().ShowBookmarks()
	case gdk.KEY_B:
		//backend.AddBookmark(backend.PANEL_ACTIVE)
	}
}

func keyBookmark(key uint) {
	b := panel.Active().Tab().Bookmarks()
	switch key {
	case gdk.KEY_j, gdk.KEY_Down:
		b.SetSelectionRel(1)
	case gdk.KEY_J, gdk.KEY_Page_Down:
		b.SetSelectionRel(20)
	case gdk.KEY_k, gdk.KEY_Up:
		b.SetSelectionRel(-1)
	case gdk.KEY_K, gdk.KEY_Page_Up:
		b.SetSelectionRel(-20)
	case gdk.KEY_g, gdk.KEY_Home:
		b.SetSelection(0)
	case gdk.KEY_G, gdk.KEY_End:
		b.SetSelection(-1)
	case gdk.KEY_b, gdk.KEY_q, gdk.KEY_Escape:
		panel.Active().Tab().HideBookmarks()
	}
}
