package gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tbuen/gocmd/internal/backend"
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
	dir := backend.GetDirectory(backend.PANEL_ACTIVE)
	switch key {
	case gdk.KEY_Q:
		// TODO: Ctrl-Q, Alt-Q etc. should not work...
		win.Close()
	case gdk.KEY_Tab:
		backend.TogglePanel()
	case gdk.KEY_r:
		dir.Reload()
	case gdk.KEY_j, gdk.KEY_Down:
		dir.SetSelectionRelative(1)
	case gdk.KEY_J, gdk.KEY_Page_Down:
		dir.SetSelectionRelative(20)
	case gdk.KEY_k, gdk.KEY_Up:
		dir.SetSelectionRelative(-1)
	case gdk.KEY_K, gdk.KEY_Page_Up:
		dir.SetSelectionRelative(-20)
	case gdk.KEY_g, gdk.KEY_Home:
		dir.SetSelectionAbsolute(0)
	case gdk.KEY_G, gdk.KEY_End:
		dir.SetSelectionAbsolute(-1)
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
	case gdk.KEY_s:
		mode = MODE_SORT
		Refresh()
	case gdk.KEY_v:
		mode = MODE_VIEW
		Refresh()
	case gdk.KEY_t:
		backend.CreateTab(backend.PANEL_ACTIVE)
	case gdk.KEY_T:
		backend.DuplicateTab(backend.PANEL_ACTIVE)
	case gdk.KEY_w:
		backend.DeleteTab(backend.PANEL_ACTIVE)
	case gdk.KEY_h, gdk.KEY_Left:
		backend.PrevTab(backend.PANEL_ACTIVE)
	case gdk.KEY_H: // TODO Shift+Left
		backend.FirstTab(backend.PANEL_ACTIVE)
	case gdk.KEY_l, gdk.KEY_Right:
		backend.NextTab(backend.PANEL_ACTIVE)
	case gdk.KEY_L: // TODO Shift+Right
		backend.LastTab(backend.PANEL_ACTIVE)
	}
}
