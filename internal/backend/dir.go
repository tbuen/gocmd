package backend

import (
	"github.com/tbuen/gocmd/internal/config"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"os/exec"
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

const (
	SORT_NAME = iota
	SORT_EXT
	SORT_SIZE
	SORT_TIME
)

type DirInfo struct {
	NumDirs, NumFiles, NumSelectedDirs, NumSelectedFiles int
	SizeFiles, SizeSelectedFiles                         int64
}

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
	View()
	Edit()
	Root()
	Home()
	Sort() (int, bool)
	SetSort(crit int, desc bool)
	Info() (info DirInfo)
}

type dir struct {
	state          int
	path           string
	ch             chan int
	files          []File
	sortCrit       int
	sortDesc       bool
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
			file := d.files[d.selection]
			if file.Dir() {
				d.dispOffsetHist[d.path] = d.dispOffset
				d.path = file.Path()
				d.dispOffset = 0
				d.selection = 0
				d.Reload()
			} else {
				cmd, args := config.FileCmd(file.Ext())
				if cmd != "" {
					args = append(args, file.Path())
					log.Println(log.DIR, "Exec command:", cmd, args)
					command := exec.Command(cmd, args...)
					err := command.Start()
					if err != nil {
						log.Println(log.DIR, "Failed: ", err)
					}
				}
			}
		}
	}
}

func (d *dir) View() {
	if d.state == STATE_IDLE {
		if d.selection < len(d.files) {
			file := d.files[d.selection]
			if !file.Dir() {
				cmd, args := config.View()
				if cmd != "" {
					args = append(args, file.Path())
					log.Println(log.DIR, "Exec command:", cmd, args)
					command := exec.Command(cmd, args...)
					err := command.Start()
					if err != nil {
						log.Println(log.DIR, "Failed: ", err)
					}
				}
			}
		}
	}
}

func (d *dir) Edit() {
	if d.state == STATE_IDLE {
		if d.selection < len(d.files) {
			file := d.files[d.selection]
			if !file.Dir() {
				cmd, args := config.Edit()
				if cmd != "" {
					args = append(args, file.Path())
					log.Println(log.DIR, "Exec command:", cmd, args)
					command := exec.Command(cmd, args...)
					err := command.Start()
					if err != nil {
						log.Println(log.DIR, "Failed: ", err)
					}
				}
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

func (d *dir) Sort() (int, bool) {
	return d.sortCrit, d.sortDesc
}

func (d *dir) SetSort(crit int, desc bool) {
	if d.sortCrit != crit || d.sortDesc != desc {
		d.sortCrit = crit
		d.sortDesc = desc
		d.sort()
		guiRefresh()
	}
}

func (d *dir) Info() (info DirInfo) {
	for _, f := range d.files {
		if f.Dir() {
			info.NumDirs++
			if f.Marked() {
				info.NumSelectedDirs++
			}
		} else {
			info.NumFiles++
			size := f.Size()
			info.SizeFiles += size
			if f.Marked() {
				info.NumSelectedFiles++
				info.SizeSelectedFiles += size
			}
		}
	}
	return
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
		if d.selection < 0 || d.selection >= len(d.files) {
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
	if d.selection < len(d.files) {
		d.files[d.selection].toggleMark()
		guiRefresh()
	}
}

func (d *dir) ToggleMarkAll() {
	for _, f := range d.files {
		f.toggleMark()
	}
	guiRefresh()
}

func (d *dir) sort() {
	if d.sortCrit == SORT_NAME {
		if d.sortDesc {
			orderedBy(dirFirst, nameDesc).sort(d.files)
		} else {
			orderedBy(dirFirst, nameAsc).sort(d.files)
		}
	} else if d.sortCrit == SORT_EXT {
		if d.sortDesc {
			orderedBy(dirFirst, extDesc, nameDesc).sort(d.files)
		} else {
			orderedBy(dirFirst, extAsc, nameAsc).sort(d.files)
		}
	} else if d.sortCrit == SORT_SIZE {
		if d.sortDesc {
			orderedBy(dirFirst, sizeDesc, nameAsc).sort(d.files)
		} else {
			orderedBy(dirFirst, sizeAsc, nameAsc).sort(d.files)
		}
	} else if d.sortCrit == SORT_TIME {
		if d.sortDesc {
			orderedBy(dirFirst, timeDesc, nameAsc).sort(d.files)
		} else {
			orderedBy(dirFirst, timeAsc, nameAsc).sort(d.files)
		}
	}
}

func reloadRoutine(d *dir) {
	for i := <-d.ch; i != 0; i = <-d.ch {
		if i == CMD_RELOAD {
			log.Println(log.DIR, "go routine for path", d.path, "received CMD_RELOAD")
			success := false
			if dir, err := os.Open(d.path); err == nil {
				if names, err := dir.Readdirnames(0); err == nil {
					d.files = d.files[0:0]
					log.Println(log.DIR, "before: len:", len(d.files), "cap:", cap(d.files))
					for _, name := range names {
						//time.Sleep(100 * time.Millisecond)
						if name[0] != '.' {
							log.Println(log.DIR, "creating file", d.path+string(filepath.Separator)+name)
							if file := newFile(d.path + string(filepath.Separator) + name); file != nil {
								d.files = append(d.files, file)
							} else {
								log.Println(log.DIR, "Failed!!")
							}
						}
					}
					d.sort()
					for i, f := range d.files {
						if f.Name() == d.selectDir {
							d.selection = i
							d.selectDir = ""
							break
						}
					}
					log.Println(log.DIR, "after: len:", len(d.files), "cap:", cap(d.files))
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
		} else {
			log.Println(log.DIR, "go routine for path", d.path, "received unknown command")
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
