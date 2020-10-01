package fs

import (
	"github.com/tbuen/gocmd/internal/config"
	"os"
	"path/filepath"
)

type File interface {
	Name() string
	Ext() string
	IsDir() bool
	IsMarked() bool
	Color() config.Color
	toggleMark()
}

type file struct {
	name      string
	ext       string
	directory bool
	regular   bool
	symlink   bool
	marked    bool
	color     config.Color
}

func newFile(info os.FileInfo) File {
	f := file{}
	f.name = info.Name()
	f.ext = filepath.Ext(f.name)
	if len(f.ext) > 0 {
		f.ext = f.ext[1:]
	}
	f.directory = info.IsDir()
	f.regular = info.Mode().IsRegular()
	f.symlink = info.Mode()&os.ModeSymlink != 0
	f.marked = false
	f.color = config.FileColor(f.ext)
	return &f
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Ext() string {
	return f.ext
}

func (f *file) IsDir() bool {
	return f.directory
}

func (f *file) IsMarked() bool {
	return f.marked
}

func (f *file) Color() config.Color {
	return f.color
}

func (f *file) toggleMark() {
	f.marked = !f.marked
}
