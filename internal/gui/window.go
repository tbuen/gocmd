package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/fs"
	"github.com/tbuen/gocmd/internal/log"
)

var window *gtk.ApplicationWindow

func init() {
	fs.RegisterRefresh(Refresh)
}

func NewWindow(app *gtk.Application, title string) {
	var err error
	window, err = gtk.ApplicationWindowNew(app)
	if err != nil {
		log.Fatalln("Could not create application window: ", err)
	}

	da, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatalln("Could not create drawing area: ", err)
	}

	window.Connect("delete-event", onDelete)
	window.Connect("key-press-event", onKeyPress)
	da.Connect("draw", onDraw)

	window.Add(da)
	window.SetTitle(title)
	window.SetDefaultSize(800, 500)
	window.Maximize()
	window.ShowAll()
}

func Refresh() {
	if window != nil {
		window.QueueDraw()
	}
}

func onDelete(win *gtk.ApplicationWindow, ev *gdk.Event) bool {
	log.Println(log.GUI, "delete")
	// return true to keep window open
	return false
}

func onKeyPress(win *gtk.ApplicationWindow, ev *gdk.Event) {
	keyEvent := gdk.EventKeyNewFromEvent(ev)
	log.Println(log.GUI, "key pressed:", keyEvent.KeyVal())
	switch keyEvent.KeyVal() {
	case gdk.KEY_q:
		// TODO: Ctrl-q, Alt-q etc. should not work...
		win.Close()
	case gdk.KEY_Tab:
		fs.TogglePanel()
	case gdk.KEY_r:
		fs.GetDirectory(fs.PANEL_ACTIVE).Reload()
	case gdk.KEY_j, gdk.KEY_Down:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(1)
	case gdk.KEY_J, gdk.KEY_Page_Down:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(20)
	case gdk.KEY_k, gdk.KEY_Up:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(-1)
	case gdk.KEY_K, gdk.KEY_Page_Up:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionRelative(-20)
	case gdk.KEY_g, gdk.KEY_Home:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionAbsolute(0)
	case gdk.KEY_G, gdk.KEY_End:
		fs.GetDirectory(fs.PANEL_ACTIVE).SetSelectionAbsolute(-1)
	case gdk.KEY_m:
		fs.GetDirectory(fs.PANEL_ACTIVE).ToggleMarkSelected()
	case gdk.KEY_M:
		fs.GetDirectory(fs.PANEL_ACTIVE).ToggleMarkAll()
	case gdk.KEY_u, gdk.KEY_numbersign:
		fs.GetDirectory(fs.PANEL_ACTIVE).GoUp()
	case gdk.KEY_Return:
		fs.GetDirectory(fs.PANEL_ACTIVE).Enter()
	case gdk.KEY_asciicircum:
		fs.GetDirectory(fs.PANEL_ACTIVE).Root()
	case gdk.KEY_asciitilde:
		fs.GetDirectory(fs.PANEL_ACTIVE).Home()
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
