package backend

import (
	"github.com/tbuen/gocmd/internal/config"
	"os"
	"path/filepath"
	"time"
)

type File interface {
	Path() string
	Name() string
	Ext() string
	Color() config.Color
	Dir() bool
	Pipe() bool
	Socket() bool
	Link() (bool, bool, string)
	//Size()
	Time() time.Time
	Size() int64
	Marked() bool
	toggleMark()
}

type file struct {
	path, name, ext         string
	color                   config.Color
	directory, pipe, socket bool
	link, linkOk            bool
	linkTarget              string
	time                    time.Time
	size                    int64
	marked                  bool
}

func newFile(path string) *file {
	path = filepath.Clean(path)
	info, err := os.Lstat(path)
	if err != nil {
		return nil
	}
	f := new(file)
	f.path = path
	f.name = info.Name()
	f.ext = filepath.Ext(f.name)
	if len(f.ext) > 0 {
		f.ext = f.ext[1:]
	}
	f.color = config.FileColor(f.ext)
	f.directory = info.IsDir()
	f.pipe = info.Mode()&os.ModeNamedPipe != 0
	f.socket = info.Mode()&os.ModeSocket != 0

	f.time = info.ModTime()
	f.size = info.Size()

	f.link = info.Mode()&os.ModeSymlink != 0
	if f.link {
		if f.linkTarget, err = os.Readlink(f.path); err == nil {
			var targetName string
			if f.linkTarget[0] == filepath.Separator {
				targetName = filepath.Clean(f.linkTarget)
			} else {
				targetName = filepath.Clean(filepath.Dir(f.path) + string(filepath.Separator) + f.linkTarget)
			}
			if target := newFile(targetName); target != nil {
				f.linkOk = !target.link || target.linkOk
				if f.linkOk {
					f.directory = target.directory
					f.pipe = target.pipe
					f.socket = target.socket
				}
			}
		}
	}
	return f
}

func (f *file) Path() string {
	return f.path
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Ext() string {
	return f.ext
}

func (f *file) Color() config.Color {
	return f.color
}

func (f *file) Dir() bool {
	return f.directory
}

func (f *file) Pipe() bool {
	return f.pipe
}

func (f *file) Socket() bool {
	return f.socket
}

func (f *file) Link() (bool, bool, string) {
	return f.link, f.linkOk, f.linkTarget
}

func (f *file) Time() time.Time {
	return f.time
}

func (f *file) Size() int64 {
	return f.size
}

func (f *file) Marked() bool {
	return f.marked
}

func (f *file) toggleMark() {
	f.marked = !f.marked
}
