package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend/gui"
	"github.com/tbuen/gocmd/internal/backend/panel"
	"github.com/tbuen/gocmd/internal/backend/tab"
	"github.com/tbuen/gocmd/internal/log"
)

var window *gtk.ApplicationWindow

func init() {
	gui.RegisterRefresh(Refresh)
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

	var x, y, w, h float64

	context.SetAntialias(cairo.ANTIALIAS_NONE)
	layout := pango.CairoCreateLayout(context)
	//layout.SetFontDescription(pango.FontDescriptionFromString("DejaVu Sans Mono 10"));
	layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro Semibold 8"))
	//layout.SetFontDescription(pango.FontDescriptionFromString("Cantarell 10"));
	context.SetLineWidth(1)

	setSourceColor(context, "C0C0C0")
	context.Rectangle(0, 0, width, height)
	context.Fill()

	if mode == MODE_SORT {
		layout.SetFontDescription(pango.FontDescriptionFromString("Source Code Pro 8"))
		context.Save()
		context.Translate(0, 0)
		drawSort(context, layout, width, 30.0, panel.Left().Tab().Directory())
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

	heightTabs := 25.0

	context.Save()
	x, y, w, h = 5, offsetTop+5, width/2-8, heightTabs
	context.Translate(x, y)
	context.Rectangle(-1, -1, w+1, h+1)
	context.Clip()
	sx1, sx2 := drawTabs(context, layout, w, h, panel.Left().Header())
	context.Restore()

	context.Save()
	x, y, w, h = 5, offsetTop+heightTabs+5, width/2-8, height-offsetTop-heightTabs-10
	context.Translate(x, y)
	context.Rectangle(-1, -1, w+1, h+1)
	context.Clip()
	switch panel.Left().Tab().Mode() {
	case tab.MODE_DIRECTORY:
		drawPanel(context, layout, w, h, sx1, sx2, panel.Left().IsActive(), panel.Left().Tab().Directory())
	case tab.MODE_BOOKMARKS:
		drawBookmarks(context, layout, w, h, sx1, sx2, panel.Left().IsActive(), panel.Left().Tab().Bookmarks())
	}
	context.Restore()

	context.Save()
	x, y, w, h = width/2+5, offsetTop+5, width/2-8, heightTabs
	context.Translate(x, y)
	context.Rectangle(-1, -1, w+1, h+1)
	context.Clip()
	sx1, sx2 = drawTabs(context, layout, w, h, panel.Right().Header())
	context.Restore()

	context.Save()
	x, y, w, h = width/2+5, offsetTop+heightTabs+5, width/2-8, height-offsetTop-heightTabs-10
	context.Translate(x, y)
	context.Rectangle(-1, -1, w+1, h+1)
	context.Clip()
	switch panel.Right().Tab().Mode() {
	case tab.MODE_DIRECTORY:
		drawPanel(context, layout, w, h, sx1, sx2, panel.Right().IsActive(), panel.Right().Tab().Directory())
	case tab.MODE_BOOKMARKS:
		drawBookmarks(context, layout, w, h, sx1, sx2, panel.Right().IsActive(), panel.Right().Tab().Bookmarks())
	}
	context.Restore()
}
