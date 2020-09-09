package fs

/*import (
	"fmt"
	"os"
)*/

type dir struct {
	Path string
}

/*func ReadDir(path string) {
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
}*/
