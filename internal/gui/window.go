package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func ShowWindow(application *gtk.Application, title string) {
	window, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Could not create application window.", err)
	}
	window.SetTitle(title)
	window.SetDefaultSize(400, 400)
	window.Show()
}
