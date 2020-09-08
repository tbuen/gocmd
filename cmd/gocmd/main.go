package main

import (
	"github.com/tbuen/gocmd/internal/app"
	"log"
	"os"
)

func main() {
	res := app.Run()
	log.Println("Application exits:", res)
	os.Exit(res)
}
