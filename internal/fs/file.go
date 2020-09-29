package fs

import (
	"github.com/tbuen/gocmd/internal/config"
	"os"
	"path/filepath"
)

type File interface {
	Name() string
	IsDir() bool
	IsMarked() bool
	Color() config.Color
	toggleMark()
}

type file struct {
	name      string
	directory bool
	regular   bool
	symlink   bool
	marked    bool
	color     config.Color
}

func newFile(info os.FileInfo) File {
	f := file{}
	f.name = info.Name()
	f.directory = info.IsDir()
	f.regular = info.Mode().IsRegular()
	f.symlink = info.Mode()&os.ModeSymlink != 0
	f.marked = false
	ext := filepath.Ext(f.name)
	if len(ext) > 0 {
		ext = ext[1:]
	}
	f.color = config.FileColor(ext)
	return &f
}

func (f *file) Name() string {
	return f.name
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
