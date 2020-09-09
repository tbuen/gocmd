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
var runIdle = true

func Run() int {
	fmt.Println(name, version)

	_, err := glib.IdleAdd(onIdle)
	if err != nil {
		log.Fatal("Could not add idle function.", err)
	}

	application, err := gtk.ApplicationNew("com.github.tbuen.gocmd", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application.", err)
	}

	application.Connect("startup", func() { onStartup(application) })
	application.Connect("activate", func() { onActivate(application) })
	application.Connect("shutdown", func() { onShutdown(application) })

	return application.Run(os.Args)
}

func onStartup(application *gtk.Application) {
	log.Println("startup")
}

func onActivate(application *gtk.Application) {
	log.Println("activate")
	gui.ShowWindow(application, name)
}

func onShutdown(application *gtk.Application) {
	log.Println("shutdown")
	runIdle = false
}

func onIdle() bool {
	//log.Println("idle")
	return runIdle
}