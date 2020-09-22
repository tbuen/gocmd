package fs

import (
	"os"
)

type File interface {
	Name() string
	IsDir() bool
	IsMarked() bool
	ToggleMark()
}

type file struct {
	name      string
	directory bool
	regular   bool
	symlink   bool
	marked    bool
}

func newFile(info os.FileInfo) File {
	f := file{}
	f.name = info.Name()
	f.directory = info.IsDir()
	f.regular = info.Mode().IsRegular()
	f.symlink = info.Mode()&os.ModeSymlink != 0
	f.marked = false
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

func (f *file) ToggleMark() {
	f.marked = !f.marked
	guiRefresh()
}
