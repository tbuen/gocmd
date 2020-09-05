package app

import "fmt"
import "github.com/tbuen/gocmd/internal/backend/directory"

var version = "devel"

func Start() {
	fmt.Println("Hello World!")
	fmt.Println("Version:", version)
	directory.ReadDir("/home/thomas")
}
