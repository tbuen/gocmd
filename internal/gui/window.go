package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend"
	"github.com/tbuen/gocmd/internal/log"
)

var window *gtk.ApplicationWindow

func init() {
	backend.RegisterRefresh(Refresh)
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

func onDraw(da *gtk.DrawingArea, context *cairo.Context) {
	log.Println(log.GUI, "draw")
	width := float64(da.GetAllocatedWidth())
	height := float64(da.GetAllocatedHeight())
	offsetTop := 0.0

	context.SetAntialias(cairo.ANTIALIAS_NONE)
	layout := pango.CairoCreateLayout(context)
	//layout.SetFontDescription(pango.FontDescriptionFromString("DejaVu Sans Mono 10"));
	layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro Semibold 8"))
	//layout.SetFontDescription(pango.FontDescriptionFromString("Cantarell 10"));
	context.SetLineWidth(1)

	if mode == MODE_SORT {
		layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro 8"))
		context.Save()
		context.Translate(0, 0)
		drawSort(context, layout, width, 30.0, backend.GetDirectory(backend.PANEL_ACTIVE))
		context.Restore()
		offsetTop = 30.0
	}
	if mode == MODE_VIEW {
		layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro 8"))
		context.Save()
		context.Translate(0, 0)
		drawView(context, layout, width, 30.0)
		context.Restore()
		offsetTop = 30.0
	}

	layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro Semibold 8"))

	context.Save()
	context.Translate(0, offsetTop)
	drawPanel(context, layout, width/2, height-offsetTop, backend.ActivePanel() == backend.PANEL_LEFT, backend.GetDirectory(backend.PANEL_LEFT))
	context.Restore()

	context.Save()
	context.Translate(width/2, offsetTop)
	drawPanel(context, layout, width/2, height-offsetTop, backend.ActivePanel() == backend.PANEL_RIGHT, backend.GetDirectory(backend.PANEL_RIGHT))
	context.Restore()
}
