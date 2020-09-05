package directory

import (
	"fmt"
	"os"
)

func ReadDir(path string) {
	fmt.Println(path)
	if dir, err := os.Open(path); err == nil {
		if names, err := dir.Readdirnames(0); err == nil {
			fmt.Println(names)
		} else {
			fmt.Println("error reading", path)
		}
		dir.Close()
	} else {
		fmt.Println("error opening", path)
	}
}
