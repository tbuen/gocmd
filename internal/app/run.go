package app

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tbuen/gocmd/internal/backend"
	"github.com/tbuen/gocmd/internal/backend/dir"
	"github.com/tbuen/gocmd/internal/config"
	"github.com/tbuen/gocmd/internal/gui"
	"github.com/tbuen/gocmd/internal/log"
	"os"
)

const name = "gocmd"

var version = "develop"
var runIdle = true

func Run() int {
	log.Println(log.MAIN, name, version)

	application, err := gtk.ApplicationNew("com.github.tbuen.gocmd", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatalln("could not create application:", err)
	}

	application.Connect("startup", func() { onStartup(application) })
	application.Connect("activate", func() { onActivate(application) })
	application.Connect("shutdown", func() { onShutdown(application) })

	return application.Run(os.Args)
}

func onStartup(application *gtk.Application) {
	log.Println(log.MAIN, "startup")
	_, err := glib.IdleAdd(onIdle)
	if err != nil {
		log.Fatalln("Could not register idle function:", err)
	}
	config.Load()
	backend.Start()
}

func onActivate(application *gtk.Application) {
	log.Println(log.MAIN, "activate")
	window := application.GetActiveWindow()
	if window == nil {
		gui.NewWindow(application, name+" "+version)
	} else {
		window.Present()
	}
}

func onShutdown(application *gtk.Application) {
	log.Println(log.MAIN, "shutdown")
	backend.Stop()
	config.Save()
	runIdle = false
}

func onIdle() bool {
	// TODO move to backend.go
	dir.Receive()
	return runIdle
}
