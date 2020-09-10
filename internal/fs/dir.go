package fs

import (
	"fmt"
	"log"
	"os"
)

const (
	STATE_IDLE   = 0
	STATE_RELOAD = 1
)

type Directory interface {
	Path() string
	reload()
}

type dir struct {
	state int
	path  string
	ch    chan int
}

func newDirectory(path string) Directory {
	return &dir{path: path}
}

func (d *dir) Path() string {
	return d.path
}

func (d *dir) reload() {
	log.Println("Reload:", d.path)
	if d.state == STATE_IDLE {
		d.state = STATE_RELOAD
		if d.ch == nil {
			log.Println("create go routine...")
			d.ch = make(chan int, 1)
			go reloadRoutine(d)
		}
		d.ch <- 5
	}
}

func reloadRoutine(d *dir) {
	i := <-d.ch
	log.Println("go routine for path", d.path, "received", i)
	if dir, err := os.Open(d.path); err == nil {
		if names, err := dir.Readdirnames(0); err == nil {
			fmt.Println(names)
		} else {
			fmt.Println("error reading", d.path)
		}
		dir.Close()
	} else {
		fmt.Println("error opening", d.path)
	}
	log.Println("go routine for path", d.path, "exiting...")
}
