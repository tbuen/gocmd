package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
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

	win.Connect("key-press-event", onKeyPress)
	da.Connect("draw", onDraw)

	win.Add(da)
	win.SetTitle(title)
	win.SetDefaultSize(400, 400)
	win.ShowAll()
}

func onDraw(da *gtk.DrawingArea, cr *cairo.Context) bool {
	cr.SetSourceRGB(0, 0, 0)
	cr.Rectangle(x*unitSize, y*unitSize, unitSize, unitSize)
	cr.Fill()
	return false
}

func onKeyPress(win *gtk.ApplicationWindow, ev *gdk.Event) bool {
	keyEvent := gdk.EventKeyNewFromEvent(ev)
	log.Println("Key:", keyEvent.KeyVal())
	if move, found := keyMap[keyEvent.KeyVal()]; found {
		move()
		win.QueueDraw()
	}
	return false
}
