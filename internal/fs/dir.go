package fs

import (
	"log"
	"os"
	"time"
)

const (
	STATE_IDLE   = 0
	STATE_RELOAD = 1
	STATE_ERROR  = 2
)

type Directory interface {
	State() int
	Path() string
	Files() []File
	Reload()
	Selection() int
	SetSelectionRelative(n int)
	SetSelectionAbsolute(n int)
	DispOffset() int
}

type dir struct {
	state      int
	path       string
	ch         chan int
	files      []File
	selection  int
	dispOffset int
}

type msg struct {
	d *dir
}

var ch = make(chan msg, 1)

func newDirectory(path string) Directory {
	return &dir{path: path}
}

func (d *dir) State() int {
	return d.state
}

func (d *dir) Path() string {
	return d.path
}

func (d *dir) Files() []File {
	return d.files
}

func (d *dir) Reload() {
	log.Println("Reload:", d.path)
	if d.state == STATE_IDLE {
		d.state = STATE_RELOAD
		if d.ch == nil {
			log.Println("create go routine...")
			d.ch = make(chan int, 1)
			go reloadRoutine(d)
		}
		d.ch <- 5
		guiRefresh()
		//close(d.ch)
	}
}

func (d *dir) Selection() int {
	return d.selection
}

func (d *dir) SetSelectionRelative(n int) {
	if n > 0 {
		d.SetSelectionAbsolute(d.selection + n)
	} else {
		n = -n
		if n > d.selection {
			n = d.selection
		}
		d.SetSelectionAbsolute(d.selection - n)
	}
}

func (d *dir) SetSelectionAbsolute(n int) {
	d.selection = n
	if d.selection >= len(d.files) {
		d.selection = len(d.files) - 1
	}
	guiRefresh()
}

func (d *dir) DispOffset() int {
	return d.dispOffset
}

func reloadRoutine(d *dir) {
	for i := <-d.ch; i != 0; i = <-d.ch {
		log.Println("go routine for path", d.path, "received", i)
		if dir, err := os.Open(d.path); err == nil {
			if names, err := dir.Readdirnames(0); err == nil {
				log.Println(names)
				d.files = d.files[0:0]
				log.Println("vorher: len:", len(d.files), "cap:", cap(d.files))
				for _, n := range names {
					time.Sleep(100 * time.Millisecond)
					d.files = append(d.files, newFile(n))
				}
				log.Println("nachher: len:", len(d.files), "cap:", cap(d.files))
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
		m.d.state = STATE_IDLE
		guiRefresh()
	case <-wait:
	}
}
