package main

import (
	"github.com/tbuen/gocmd/internal/app"
	"github.com/tbuen/gocmd/internal/log"
	"os"
)

func main() {
	res := app.Run()
	log.Println(log.MAIN, "Application exits: ", res)
	os.Exit(res)
}
