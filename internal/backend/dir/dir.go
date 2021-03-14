package dir

import (
	"github.com/tbuen/gocmd/internal/backend/gui"
	"github.com/tbuen/gocmd/internal/config"
	. "github.com/tbuen/gocmd/internal/global"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	STATE_INIT = iota
	STATE_IDLE
	STATE_RELOAD
	STATE_ERROR
	STATE_DEAD
)

const (
	_ = iota
	CMD_RELOAD
)

type DirInfo struct {
	NumDirs, NumFiles, NumSelectedDirs, NumSelectedFiles int
	SizeFiles, SizeSelectedFiles                         int64
}

type Directory struct {
	state          int
	config         config.Directory
	ch             chan int
	files          []File
	sortKey        int
	sortOrder      int
	hidden         bool
	selection      int
	dispOffset     int
	dispOffsetHist map[string]int
	selectDir      string
}

type msg struct {
	success bool
	dir     *Directory
}

var channel = make(chan msg, 1)

func New() *Directory {
	home, err := os.UserHomeDir()
	var path string
	if err == nil {
		path = home
	} else {
		path = string(filepath.Separator)
	}
	return newDirectory(config.Directory{path}, SORT_BY_NAME, SORT_ASCENDING, false)
}

func NewWithConfig(cfg config.Directory) *Directory {
	return newDirectory(cfg, SORT_BY_NAME, SORT_ASCENDING, false)
}

func (dir *Directory) Clone() *Directory {
	return newDirectory(dir.config, dir.sortKey, dir.sortOrder, dir.hidden)
}

func (d *Directory) Config() config.Directory {
	return d.config
}

func newDirectory(cfg config.Directory, sortKey int, sortOrder int, hidden bool) *Directory {
	d := new(Directory)
	d.dispOffsetHist = make(map[string]int)
	d.config = cfg
	d.sortKey = sortKey
	d.sortOrder = sortOrder
	d.hidden = hidden
	return d
}

func (dir *Directory) State() int {
	return dir.state
}

func (dir *Directory) Path() string {
	return dir.config.Path
}

func (dir *Directory) Files() []File {
	return dir.files
}

func (dir *Directory) Reload() {
	log.Println(log.DIR, "Reload:", dir.config.Path)
	if dir.state == STATE_INIT || dir.state == STATE_IDLE || dir.state == STATE_ERROR {
		dir.state = STATE_RELOAD
		if dir.ch == nil {
			log.Println(log.DIR, "create go routine...")
			dir.ch = make(chan int, 1)
			go reloadRoutine(dir)
		}
		dir.ch <- CMD_RELOAD
		gui.Refresh()
	}
}

func (dir *Directory) Destroy() {
	if dir.state != STATE_DEAD {
		dir.state = STATE_DEAD
		if dir.ch != nil {
			close(dir.ch)
		}
	}
}

func (dir *Directory) GoUp() {
	if dir.state == STATE_IDLE || dir.state == STATE_ERROR {
		if dir.config.Path != string(filepath.Separator) {
			dir.selectDir = filepath.Base(dir.config.Path)
			dir.config.Path = filepath.Dir(dir.config.Path)
			dir.dispOffset = 0
			dir.selection = 0
			dir.Reload()
		}
	}
}

func (dir *Directory) Enter() {
	if dir.state == STATE_IDLE {
		if dir.selection < len(dir.files) {
			file := dir.files[dir.selection]
			if file.Dir() {
				dir.dispOffsetHist[dir.config.Path] = dir.dispOffset
				dir.config.Path = file.Path()
				dir.dispOffset = 0
				dir.selection = 0
				dir.Reload()
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

func (dir *Directory) View() {
	if dir.state == STATE_IDLE {
		if dir.selection < len(dir.files) {
			file := dir.files[dir.selection]
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

func (dir *Directory) Edit() {
	if dir.state == STATE_IDLE {
		if dir.selection < len(dir.files) {
			file := dir.files[dir.selection]
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

func (dir *Directory) Root() {
	if dir.state == STATE_IDLE || dir.state == STATE_ERROR {
		if dir.config.Path != string(filepath.Separator) {
			dir.config.Path = string(filepath.Separator)
			dir.dispOffset = 0
			dir.selection = 0
			dir.Reload()
		}
	}
}

func (dir *Directory) Home() {
	if dir.state == STATE_IDLE || dir.state == STATE_ERROR {
		home, err := os.UserHomeDir()
		if err == nil {
			if dir.config.Path != home {
				dir.config.Path = home
				dir.dispOffset = 0
				dir.selection = 0
				dir.Reload()
			}
		}
	}
}

func (dir *Directory) SortKey() (sortKey, sortOrder int) {
	// TODO state!!
	sortKey, sortOrder = dir.sortKey, dir.sortOrder
	return
}

func (dir *Directory) SetSortKey(sortKey, sortOrder int) {
	// TODO state!!
	if dir.sortKey != sortKey || dir.sortOrder != sortOrder {
		dir.sortKey = sortKey
		dir.sortOrder = sortOrder
		dir.sort()
		gui.Refresh()
	}
}

func (dir *Directory) Hidden() bool {
	return dir.hidden
}

func (dir *Directory) ToggleHidden() {
	if dir.state == STATE_IDLE || dir.state == STATE_ERROR {
		dir.hidden = !dir.hidden
		dir.Reload()
	}
}

func (dir *Directory) Info() (info DirInfo) {
	for _, f := range dir.files {
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

func (dir *Directory) Selection() int {
	return dir.selection
}

func (dir *Directory) SetSelectionRelative(n int) {
	if dir.state == STATE_IDLE {
		if n > 0 {
			dir.SetSelectionAbsolute(dir.selection + n)
		} else {
			n = -n
			if n > dir.selection {
				n = dir.selection
			}
			dir.SetSelectionAbsolute(dir.selection - n)
		}
	}
}

func (dir *Directory) SetSelectionAbsolute(n int) {
	if dir.state == STATE_IDLE {
		dir.selection = n
		if dir.selection < 0 || dir.selection >= len(dir.files) {
			dir.selection = len(dir.files) - 1
		}
		gui.Refresh()
	}
}

func (dir *Directory) DispOffset() int {
	return dir.dispOffset
}

func (dir *Directory) SetDispOffset(offset int) {
	dir.dispOffset = offset
}

func (dir *Directory) ToggleMarkSelected() {
	if dir.selection < len(dir.files) {
		dir.files[dir.selection].toggleMark()
		gui.Refresh()
	}
}

func (dir *Directory) ToggleMarkAll() {
	if len(dir.files) > 0 {
		for _, f := range dir.files {
			f.toggleMark()
		}
		gui.Refresh()
	}
}

func (dir *Directory) sort() {
	if dir.sortKey == SORT_BY_NAME {
		if dir.sortOrder == SORT_ASCENDING {
			orderedBy(dirFirst, nameAsc).sort(dir.files)
		} else {
			orderedBy(dirFirst, nameDesc).sort(dir.files)
		}
	} else if dir.sortKey == SORT_BY_EXT {
		if dir.sortOrder == SORT_ASCENDING {
			orderedBy(dirFirst, extAsc, nameAsc).sort(dir.files)
		} else {
			orderedBy(dirFirst, extDesc, nameDesc).sort(dir.files)
		}
	} else if dir.sortKey == SORT_BY_SIZE {
		if dir.sortOrder == SORT_ASCENDING {
			orderedBy(dirFirst, sizeAsc, nameAsc).sort(dir.files)
		} else {
			orderedBy(dirFirst, sizeDesc, nameAsc).sort(dir.files)
		}
	} else if dir.sortKey == SORT_BY_TIME {
		if dir.sortOrder == SORT_ASCENDING {
			orderedBy(dirFirst, timeAsc, nameAsc).sort(dir.files)
		} else {
			orderedBy(dirFirst, timeDesc, nameAsc).sort(dir.files)
		}
	}
}

func reloadRoutine(dir *Directory) {
	for cmd := <-dir.ch; cmd != 0; cmd = <-dir.ch {
		if cmd == CMD_RELOAD {
			log.Println(log.DIR, "go routine for path", dir.config.Path, "received CMD_RELOAD")
			success := false
			var prevSelectedFile string
			if dir.selectDir == "" && dir.selection < len(dir.files) {
				prevSelectedFile = dir.files[dir.selection].Name()
			}
			if directory, err := os.Open(dir.config.Path); err == nil {
				if names, err := directory.Readdirnames(0); err == nil {
					dir.files = dir.files[0:0]
					log.Println(log.DIR, "before: len:", len(dir.files), "cap:", cap(dir.files))
					for _, name := range names {
						//time.Sleep(100 * time.Millisecond)
						if dir.hidden || name[0] != '.' {
							//log.Println(log.DIR, "creating file", dir.config.Path+string(filepath.Separator)+name)
							if file := newFile(dir.config.Path + string(filepath.Separator) + name); file != nil {
								dir.files = append(dir.files, file)
							} else {
								log.Println(log.DIR, "Failed!!")
							}
						}
					}
					dir.sort()
					for i, f := range dir.files {
						if name := f.Name(); name == dir.selectDir || name == prevSelectedFile {
							dir.selection = i
							dir.selectDir = ""
							break
						}
					}
					if !dir.hidden && prevSelectedFile != "" && prevSelectedFile[0] == '.' {
						dir.selection = 0
					}
					log.Println(log.DIR, "after: len:", len(dir.files), "cap:", cap(dir.files))
					if offset, ok := dir.dispOffsetHist[dir.config.Path]; ok {
						dir.dispOffset = offset
						delete(dir.dispOffsetHist, dir.config.Path)
					}
					success = true
				} else {
					log.Println(log.DIR, "error reading", dir.config.Path)
				}
				directory.Close()
			} else {
				log.Println(log.DIR, "error opening", dir.config.Path)
			}
			m := msg{success, dir}
			channel <- m
		} else {
			log.Println(log.DIR, "go routine for path", dir.config.Path, "received unknown command")
		}
	}
	log.Println(log.DIR, "go routine for path", dir.config.Path, "exiting...")
}

func Receive() {
	wait := time.After(10 * time.Millisecond)
	select {
	case msg := <-channel:
		log.Println(log.DIR, "received response for path", msg.dir.Path())
		if msg.dir.state == STATE_RELOAD {
			if msg.success {
				msg.dir.state = STATE_IDLE
			} else {
				msg.dir.state = STATE_ERROR
			}
			gui.Refresh()
		} else {
			log.Println(log.DIR, "directory already dead")
		}
	case <-wait:
	}
}