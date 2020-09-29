package fs

import (
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"path/filepath"
	"time"
)

const (
	STATE_IDLE = iota
	STATE_RELOAD
	STATE_ERROR
)

const (
	CMD_RELOAD = 1
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
	SetDispOffset(offset int)
	ToggleMarkSelected()
	ToggleMarkAll()
	GoUp()
	Enter()
	Root()
	Home()
}

type dir struct {
	state          int
	path           string
	ch             chan int
	files          []File
	selection      int
	dispOffset     int
	dispOffsetHist map[string]int
	selectDir      string
}

type msg struct {
	success bool
	d       *dir
}

var ch = make(chan msg, 1)

func newDirectory(path string) Directory {
	d := dir{}
	d.state = STATE_ERROR
	d.dispOffsetHist = make(map[string]int)
	if path == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			d.path = home
		} else {
			d.path = string(filepath.Separator)
		}
	} else {
		d.path = path
	}
	return &d
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
	log.Println(log.DIR, "Reload:", d.path)
	if d.state != STATE_RELOAD {
		d.state = STATE_RELOAD
		if d.ch == nil {
			log.Println(log.DIR, "create go routine...")
			d.ch = make(chan int, 1)
			go reloadRoutine(d)
		}
		d.ch <- CMD_RELOAD
		guiRefresh()
		//close(d.ch)
	}
}

func (d *dir) GoUp() {
	if d.state != STATE_RELOAD {
		if d.path != string(filepath.Separator) {
			d.selectDir = filepath.Base(d.path)
			d.path = filepath.Dir(d.path)
			d.dispOffset = 0
			d.selection = 0
			d.Reload()
		}
	}
}

func (d *dir) Enter() {
	if d.state == STATE_IDLE {
		if d.selection < len(d.files) {
			if d.files[d.selection].IsDir() {
				d.dispOffsetHist[d.path] = d.dispOffset
				if d.path == string(filepath.Separator) {
					d.path += d.files[d.selection].Name()
				} else {
					d.path += string(filepath.Separator) + d.files[d.selection].Name()
				}
				d.dispOffset = 0
				d.selection = 0
				d.Reload()
			} else {
				// TODO
			}
		}
	}
}

func (d *dir) Root() {
	if d.state != STATE_RELOAD {
		if d.path != string(filepath.Separator) {
			d.path = string(filepath.Separator)
			d.dispOffset = 0
			d.selection = 0
			d.Reload()
		}
	}
}

func (d *dir) Home() {
	if d.state != STATE_RELOAD {
		home, err := os.UserHomeDir()
		if err == nil {
			if d.path != home {
				d.path = home
				d.dispOffset = 0
				d.selection = 0
				d.Reload()
			}
		}
	}
}

func (d *dir) Selection() int {
	return d.selection
}

func (d *dir) SetSelectionRelative(n int) {
	if d.state == STATE_IDLE {
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
}

func (d *dir) SetSelectionAbsolute(n int) {
	if d.state == STATE_IDLE {
		d.selection = n
		if d.selection >= len(d.files) {
			d.selection = len(d.files) - 1
		}
		guiRefresh()
	}
}

func (d *dir) DispOffset() int {
	return d.dispOffset
}

func (d *dir) SetDispOffset(offset int) {
	d.dispOffset = offset
}

func (d *dir) ToggleMarkSelected() {
	d.files[d.selection].toggleMark()
	guiRefresh()
}

func (d *dir) ToggleMarkAll() {
	for _, f := range d.files {
		f.toggleMark()
	}
	guiRefresh()
}

func reloadRoutine(d *dir) {
	for i := <-d.ch; i != 0; i = <-d.ch {
		log.Println(log.DIR, "go routine for path", d.path, "received", i)
		if i == CMD_RELOAD {
			success := false
			if dir, err := os.Open(d.path); err == nil {
				if fileinfo, err := dir.Readdir(0); err == nil {
					d.files = d.files[0:0]
					log.Println(log.DIR, "vorher: len:", len(d.files), "cap:", cap(d.files))
					for _, fi := range fileinfo {
						log.Println(log.DIR, "Datei: ", fi.Name())
						//time.Sleep(100 * time.Millisecond)
						if fi.Name()[0] != '.' {
							d.files = append(d.files, newFile(fi))
						}
					}
					orderedBy(dirFirst, name).sort(d.files)
					for i, f := range d.files {
						if f.Name() == d.selectDir {
							d.selection = i
							d.selectDir = ""
							break
						}
					}
					log.Println(log.DIR, "nachher: len:", len(d.files), "cap:", cap(d.files))
					if offset, ok := d.dispOffsetHist[d.path]; ok {
						d.dispOffset = offset
						delete(d.dispOffsetHist, d.path)
					}
					success = true

				} else {
					log.Println(log.DIR, "error reading", d.path)
				}
				dir.Close()
			} else {
				log.Println(log.DIR, "error opening", d.path)

			}
			m := msg{success, d}
			ch <- m
		}
	}
	log.Println(log.DIR, "go routine for path", d.path, "exiting...")
}

func Receive() {
	wait := time.After(10 * time.Millisecond)
	select {
	case m := <-ch:
		log.Println(log.DIR, "received response for path", m.d.Path())
		if m.success {
			m.d.state = STATE_IDLE
		} else {
			m.d.state = STATE_ERROR
		}
		guiRefresh()
	case <-wait:
	}
}
