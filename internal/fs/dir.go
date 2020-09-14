package fs

import (
	"log"
	"os"
	"time"
)

const (
	STATE_IDLE   = 0
	STATE_RELOAD = 1
)

type Directory interface {
	Path() string
	Files() []File
	reload()
}

type dir struct {
	state int
	path  string
	ch    chan int
	files []File
}

type msg struct {
	d *dir
}

var ch = make(chan msg, 1)

func newDirectory(path string) Directory {
	return &dir{path: path}
}

func (d *dir) Path() string {
	return d.path
}

func (d *dir) Files() []File {
	return d.files
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
		//close(d.ch)
	}
}

func reloadRoutine(d *dir) {
	for i := <-d.ch; i != 0; i = <-d.ch {
		log.Println("go routine for path", d.path, "received", i)
		if dir, err := os.Open(d.path); err == nil {
			if names, err := dir.Readdirnames(0); err == nil {
				log.Println(names)
				for _, n := range names {
					d.files = append(d.files, newFile(n))
				}
			} else {
				log.Println("error reading", d.path)
			}
			dir.Close()
		} else {
			log.Println("error opening", d.path)
		}
		m := msg{d}
		ch <- m
	}
	log.Println("go routine for path", d.path, "exiting...")
}

func Receive() {
	wait := time.After(10 * time.Millisecond)
	select {
	case m := <-ch:
		log.Println("received response for path", m.d.Path())
	case <-wait:
	}
}
