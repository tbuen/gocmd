package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tbuen/gocmd/internal/fs"
	"log"
)

const (
	KEY_LEFT  uint = 65361
	KEY_UP    uint = 65362
	KEY_RIGHT uint = 65363
	KEY_DOWN  uint = 65364
)

var (
	unitSize = 20.0
	x        = 0.0
	y        = 0.0
	keyMap   = map[uint]func(){
		KEY_LEFT:  func() { x-- },
		KEY_UP:    func() { y-- },
		KEY_RIGHT: func() { x++ },
		KEY_DOWN:  func() { y++ },
	}
	application *gtk.Application
)

func ShowWindow(app *gtk.Application, title string) {
	application = app

	win, err := gtk.ApplicationWindowNew(application)
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
	if keyEvent.KeyVal() == gdk.KEY_q {
		// TODO: Ctrl-q, Alt-q etc. should not work...
		win.Close()
	} else if move, found := keyMap[keyEvent.KeyVal()]; found {
		move()
		win.QueueDraw()
	}
}

func onDraw(da *gtk.DrawingArea, cr *cairo.Context) {
	dirLeft := fs.Tab(fs.TAB_LEFT)
	dirRight := fs.Tab(fs.TAB_RIGHT)
	log.Println("path(left) =", dirLeft.Path)
	log.Println("path(right) =", dirRight.Path)
	cr.SetSourceRGB(0, 0, 0)
	cr.Rectangle(x*unitSize, y*unitSize, unitSize, unitSize)
	cr.Fill()
}
