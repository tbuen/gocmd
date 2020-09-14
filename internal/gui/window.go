package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/fs"
	"log"
)

func ShowWindow(app *gtk.Application, title string) {
	win, err := gtk.ApplicationWindowNew(app)
	if err != nil {
		log.Fatal("Could not create application window.", err)
	}

	da, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Could not create drawing area.", err)
	}

	win.Connect("delete-event", onDelete)
	win.Connect("key-press-event", onKeyPress)
	da.Connect("draw", onDraw)

	win.Add(da)
	win.SetTitle(title)
	win.SetDefaultSize(400, 400)
	win.ShowAll()
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
	case gdk.KEY_r:
		win.QueueDraw()
	}
}

func onDraw(da *gtk.DrawingArea, context *cairo.Context) {
	width := da.GetAllocatedWidth()
	height := da.GetAllocatedHeight()

	context.SetAntialias(cairo.ANTIALIAS_NONE)
	layout := pango.CairoCreateLayout(context)
	//layout.SetFontDescription(pango.FontDescriptionFromString("DejaVu Sans Mono 10"));
	layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro Semibold 10"))
	//layout.SetFontDescription(pango.FontDescriptionFromString("Cantarell 10"));
	context.SetLineWidth(1)

	context.Save()
	context.Translate(0, 0)
	drawPanel(context, layout, float64(width/2), float64(height), fs.Tab(fs.TAB_LEFT))
	context.Restore()

	context.Save()
	context.Translate(float64(width/2), 0)
	drawPanel(context, layout, float64(width/2), float64(height), fs.Tab(fs.TAB_RIGHT))
	context.Restore()
}
