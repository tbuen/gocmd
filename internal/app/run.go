package app

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tbuen/gocmd/internal/gui"
	"log"
	"os"
)

const name = "gocmd"
var version = "develop"

func Run() {
	fmt.Println(name, version)

	application, err := gtk.ApplicationNew("com.github.tbuen.gocmd", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application.", err)
	}

	application.Connect("startup", func() { onStartup(application) })
	application.Connect("activate", func() { onActivate(application) })
	application.Connect("shutdown", func() { onShutdown(application) })

	os.Exit(application.Run(os.Args))
}

func onStartup(application *gtk.Application) {
	log.Print("startup")
}

func onActivate(application *gtk.Application) {
	log.Print("activate")
	gui.ShowWindow(application, name)
}

func onShutdown(application *gtk.Application) {
	log.Print("shutdown")
}
