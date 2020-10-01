package backend

import (
	"github.com/tbuen/gocmd/internal/config"
	"os"
	"path/filepath"
)

type File interface {
	Path() string
	Name() string
	Ext() string
	Color() config.Color
	Dir() bool
	Link() (bool, bool, string)
	Marked() bool
	toggleMark()
}

type file struct {
	path       string
	name       string
	ext        string
	color      config.Color
	directory  bool
	link       bool
	linkOk     bool
	linkTarget string
	marked     bool
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
	f.link = info.Mode()&os.ModeSymlink != 0
	if f.link {
		if f.linkTarget, err = os.Readlink(f.path); err == nil {
			if target := newFile(filepath.Clean(filepath.Dir(f.path) + string(filepath.Separator) + f.linkTarget)); target != nil {
				f.linkOk = !target.link || target.linkOk
				if f.linkOk && target.directory {
					f.directory = true
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

func (f *file) Dir() bool {
	return f.directory
}

func (f *file) Link() (bool, bool, string) {
	return f.link, f.linkOk, f.linkTarget
}

func (f *file) Color() config.Color {
	return f.color
}

func (f *file) Marked() bool {
	return f.marked
}

func (f *file) toggleMark() {
	f.marked = !f.marked
}
