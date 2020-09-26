package app

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tbuen/gocmd/internal/fs"
	"github.com/tbuen/gocmd/internal/gui"
	"github.com/tbuen/gocmd/internal/log"
	"os"
)

const name = "gocmd"

var version = "develop"
var runIdle = true

func Run() int {
	log.Println(log.MOD_MAIN, name, version)

	_, err := glib.IdleAdd(onIdle)
	if err != nil {
		log.Fatal("could not add idle function:", err)
	}

	application, err := gtk.ApplicationNew("com.github.tbuen.gocmd", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("could not create application:", err)
	}

	application.Connect("startup", func() { onStartup(application) })
	application.Connect("activate", func() { onActivate(application) })
	application.Connect("shutdown", func() { onShutdown(application) })

	return application.Run(os.Args)
}

func onStartup(application *gtk.Application) {
	log.Println(log.MOD_MAIN, "startup")
}

func onActivate(application *gtk.Application) {
	log.Println(log.MOD_MAIN, "activate")
	gui.ShowWindow(application, name+" "+version)
}

func onShutdown(application *gtk.Application) {
	log.Println(log.MOD_MAIN, "shutdown")
	runIdle = false
}

func onIdle() bool {
	fs.Receive()
	return runIdle
}
