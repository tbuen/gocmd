package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/fs"
	"log"
)

var window *gtk.ApplicationWindow

func init() {
	fs.RegisterRefresh(Refresh)
}

func ShowWindow(app *gtk.Application, title string) {
	var err error
	window, err = gtk.ApplicationWindowNew(app)
	if err != nil {
		log.Fatal("Could not create application window.", err)
	}

	da, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Could not create drawing area.", err)
	}

	window.Connect("delete-event", onDelete)
	window.Connect("key-press-event", onKeyPress)
	da.Connect("draw", onDraw)

	window.Add(da)
	window.SetTitle(title)
	window.SetDefaultSize(1000, 600)
	window.ShowAll()
}

func Refresh() {
	if window != nil {
		window.QueueDraw()
	}
}

func onDelete(win *gtk.ApplicationWindow, ev *gdk.Event) bool {
	log.Println("delete")
	// return true to keep window open
	return false
}

func onKeyPress(win *gtk.ApplicationWindow, ev *gdk.Event) {
	keyEvent := gdk.EventKeyNewFromEvent(ev)
	log.Println("Key:", keyEvent.KeyVal())
	switch keyEvent.KeyVal() {
	case gdk.KEY_q:
		// TODO: Ctrl-q, Alt-q etc. should not work...
		win.Close()
	case gdk.KEY_Tab:
		fs.TogglePanel()
	case gdk.KEY_r:
		fs.GetDirectory(fs.PANEL_ACTIVE).Reload()
	case gdk.KEY_j:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(1)
	case gdk.KEY_J:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(20)
	case gdk.KEY_k:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(-1)
	case gdk.KEY_K:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(-20)
	case gdk.KEY_g:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionAbsolute(0)
	case gdk.KEY_G:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionAbsolute(len(fs.GetDirectory(fs.PANEL_ACTIVE).Files()) - 1)
	case gdk.KEY_n:
		fs.GetDirectory(fs.PANEL_ACTIVE).Files()[fs.GetDirectory(fs.PANEL_ACTIVE).Selection()].ToggleMark()
	case gdk.KEY_N:
		for _, f := range fs.GetDirectory(fs.PANEL_ACTIVE).Files() {
			f.ToggleMark()
		}
	}
}

func onDraw(da *gtk.DrawingArea, context *cairo.Context) {
	width := da.GetAllocatedWidth()
	height := da.GetAllocatedHeight()

	context.SetAntialias(cairo.ANTIALIAS_NONE)
	layout := pango.CairoCreateLayout(context)
	//layout.SetFontDescription(pango.FontDescriptionFromString("DejaVu Sans Mono 10"));
	layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro Semibold 8"))
	//layout.SetFontDescription(pango.FontDescriptionFromString("Cantarell 10"));
	context.SetLineWidth(1)

	context.Save()
	context.Translate(0, 0)
	drawPanel(context, layout, float64(width/2), float64(height), fs.ActivePanel() == fs.PANEL_LEFT, fs.GetDirectory(fs.PANEL_LEFT))
	context.Restore()

	context.Save()
	context.Translate(float64(width/2), 0)
	drawPanel(context, layout, float64(width/2), float64(height), fs.ActivePanel() == fs.PANEL_RIGHT, fs.GetDirectory(fs.PANEL_RIGHT))
	context.Restore()
}
